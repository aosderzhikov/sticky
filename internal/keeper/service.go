package keeper

import (
	"log/slog"
	"sync"
	"time"
)

func NewService() *Keeper {
	return &Keeper{
		values:   make(map[string]value),
		ttlQueue: make(map[string]<-chan time.Time),
	}
}

type Keeper struct {
	muV    sync.Mutex
	values map[string]value

	muT      sync.Mutex
	ttlQueue map[string]<-chan time.Time
}

type value struct {
	key  string
	data []byte
	ttl  time.Duration
}

func (k *Keeper) Get(key string) []byte {
	k.muV.Lock()
	defer k.muV.Unlock()
	return k.values[key].data
}

func (k *Keeper) Set(key string, data []byte, ttl time.Duration) {
	k.muV.Lock()
	k.values[key] = value{key, data, ttl}
	k.muV.Unlock()

	k.muT.Lock()
	timer := time.NewTimer(ttl)
	k.ttlQueue[key] = timer.C
	k.muT.Unlock()
	slog.Info("setted value: %v %v", k.values, k.ttlQueue)
}

func (k *Keeper) Delete(key string) {
	k.muV.Lock()
	delete(k.values, key)
	k.muV.Unlock()

	k.muT.Lock()
	delete(k.ttlQueue, key)
	k.muT.Unlock()
}

func (k *Keeper) Run() {
	go k.observeTTL()
}

func (k *Keeper) observeTTL() {
	for {
		for key := range k.ttlQueue {
			k.muT.Lock()
			ch, ok := k.ttlQueue[key]
			if !ok {
				continue
			}
			k.muT.Unlock()

			select {
			case <-ch:
				k.Delete(key)
			default:
			}
		}
	}
}
