package keeper

import (
	"fmt"
	"log/slog"
	"maps"
	"sync"
	"time"
)

func NewService(ttl time.Duration) *Keeper {
	return &Keeper{
		values:     make(map[string]value),
		defaultTTL: ttl,
	}
}

type Keeper struct {
	mu         sync.RWMutex
	values     map[string]value
	defaultTTL time.Duration
}

type value struct {
	key  string
	data []byte
	ttl  *time.Timer
}

func (k *Keeper) Get(key string) []byte {
	k.mu.RLock()
	data := k.values[key].data
	k.mu.RUnlock()
	return data
}

func (k *Keeper) Set(key string, data []byte, ttl time.Duration) {
	k.mu.Lock()
	if ttl == 0 {
		ttl = k.defaultTTL
	}
	k.values[key] = value{key, data, time.NewTimer(ttl)}
	k.mu.Unlock()

	slog.Debug(fmt.Sprintf("set key %q with ttl %s", key, ttl))
}

func (k *Keeper) Delete(key string) {
	k.mu.RLock()
	val, ok := k.values[key]
	k.mu.RUnlock()
	if !ok {
		return
	}

	val.ttl.Stop()
	k.mu.Lock()
	delete(k.values, key)
	k.mu.Unlock()

	slog.Debug(fmt.Sprintf("delete key %q", key))
}

func (k *Keeper) Run() {
	go k.observeTTL()
}

// TODO: think about optimization
func (k *Keeper) observeTTL() {
	copied := make(map[string]value)
	for {
		k.mu.RLock()
		maps.Copy(copied, k.values)
		k.mu.RUnlock()

		for key := range copied {
			val, ok := copied[key]
			if !ok {
				continue
			}

			select {
			case <-val.ttl.C:
				slog.Debug(fmt.Sprintf("key %q expired", key))
				k.Delete(key)
			default:
			}
		}

		copied = make(map[string]value)
	}
}
