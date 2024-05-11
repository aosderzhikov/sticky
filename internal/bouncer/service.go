//go:generate mockgen -source=$GOFILE -destination=mock_test.go -package=$GOPACKAGE
package bouncer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

func NewShardService(storages []Storage) *ShardService {
	return &ShardService{
		storages: storages,
		index:    make(map[string]int),
	}
}

type ShardService struct {
	storages []Storage

	mu    sync.Mutex
	index map[string]int

	defaultTTL time.Duration
}

type Storage interface {
	Get(ctx context.Context, key string) (value []byte, err error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) (err error)
	Delete(ctx context.Context, key string) (err error)

	Addr() (addr string)
	IsAlive() (alive bool)
}

var (
	ErrAllStorage  error = errors.New("not found storage to store value")
	ErrKeyNotExist error = errors.New("key not exist")
)

func (b *ShardService) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	var s Storage

	if ttl == 0 {
		ttl = b.defaultTTL
	}

	i, exist := b.getStorageIndex(key)
	if exist && b.storages[i].IsAlive() {
		s = b.storages[i]
		err := s.Set(ctx, key, value, ttl)
		if err == nil {
			return nil
		}
	}

	// TODO: this approach leaves unused values in unavailable shards
	// if one shard is dead all updates will set in first shard
	// think about round robin or hash func
	for i := range b.storages {
		s = b.storages[i]
		if !s.IsAlive() {
			continue
		}

		err := s.Set(ctx, key, value, ttl)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("store key %q in storage with addr %q failed: %v", key, s.Addr(), err))
			continue
		}

		b.setStorageIndex(key, i)
		return nil
	}

	return ErrAllStorage
}

func (b *ShardService) Delete(ctx context.Context, key string) error {
	i, ok := b.getStorageIndex(key)
	if !ok {
		return ErrKeyNotExist
	}

	s := b.storages[i]
	if !s.IsAlive() {
		return fmt.Errorf("storage %q isnt alive", s.Addr())
	}

	return s.Delete(ctx, key)
}

func (b *ShardService) Get(ctx context.Context, key string) ([]byte, error) {
	i, ok := b.getStorageIndex(key)
	if !ok {
		return nil, ErrKeyNotExist
	}

	s := b.storages[i]
	if !s.IsAlive() {
		return nil, fmt.Errorf("storage %q isnt alive", s.Addr())
	}

	return s.Get(ctx, key)
}

func (b *ShardService) getStorageIndex(key string) (int, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	i, ok := b.index[key]
	return i, ok
}

func (b *ShardService) setStorageIndex(key string, i int) {
	b.mu.Lock()
	b.index[key] = i
	b.mu.Unlock()
}
