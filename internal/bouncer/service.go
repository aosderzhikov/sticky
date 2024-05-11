//go:generate mockgen -source=$GOFILE -destination=mock_test.go -package=$GOPACKAGE
package bouncer

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
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

	i, exist := b.isExist(key)
	if exist && b.storages[i].IsAlive() {
		slog.Debug(fmt.Sprintf("key %q is exist, value will be updated", key))
		s = b.storages[i]
		err := s.Set(ctx, key, value, ttl)
		if err == nil {
			b.setStorageIndex(key, i)
			return nil
		}
		slog.ErrorContext(ctx, fmt.Sprintf("update value by key %q failed: %v", key, err))
	}

	i = b.getShardIndByHash(key)
	if b.storages[i].IsAlive() {
		s = b.storages[i]
		slog.Debug(fmt.Sprintf("selected by hash storage with index %d and addr %q is alive", i, s.Addr()))
		err := s.Set(ctx, key, value, ttl)
		if err == nil {
			b.setStorageIndex(key, i)
			return nil
		}
		slog.ErrorContext(ctx, fmt.Sprintf("store value by key %q failed: %v", key, err))
	}

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
	i, ok := b.isExist(key)
	if !ok {
		return ErrKeyNotExist
	}

	s := b.storages[i]
	if !s.IsAlive() {
		return fmt.Errorf("storage %q isnt alive", s.Addr())
	}

	err := s.Delete(ctx, key)
	if err != nil {
		return err
	}

	b.deletStorageIndex(key)
	return nil
}

func (b *ShardService) Get(ctx context.Context, key string) ([]byte, error) {
	i, ok := b.isExist(key)
	if !ok {
		return nil, ErrKeyNotExist
	}

	s := b.storages[i]
	if !s.IsAlive() {
		return nil, fmt.Errorf("storage %q isnt alive", s.Addr())
	}

	value, err := s.Get(ctx, key)
	if len(value) == 0 {
		b.deletStorageIndex(key)
		return nil, ErrKeyNotExist
	}
	return value, err
}

func (b *ShardService) getShardIndByHash(key string) int {
	aliveShards := b.countAliveShards()
	h := fnv.New64a()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(aliveShards))
}

func (b *ShardService) countAliveShards() int {
	count := 0
	for _, s := range b.storages {
		if s.IsAlive() {
			count++
		}
	}
	return count
}

func (b *ShardService) isExist(key string) (int, bool) {
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

func (b *ShardService) deletStorageIndex(key string) {
	b.mu.Lock()
	delete(b.index, key)
	b.mu.Unlock()
}
