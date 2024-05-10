package bouncer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func NewShard(addr string, interval time.Duration, client *http.Client) *Shard {
	if client == nil {
		client = &http.Client{}
	}

	if !strings.HasSuffix(addr, "/") {
		addr += "/"
	}

	return &Shard{
		addr:                addr,
		client:              *client,
		healthCheckInterval: interval,
	}
}

type Shard struct {
	addr                string
	alive               bool
	healthCheckInterval time.Duration
	client              http.Client
}

const (
	setEndpoint         = "set"
	getEndpoint         = "get"
	deleteEndpoint      = "delete"
	healthCheckEndpoint = "health-check"
)

func (s *Shard) Get(ctx context.Context, key string) (value []byte, err error) {
	url := s.addr + getEndpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	putKey(req, key)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *Shard) Set(ctx context.Context, key string, value []byte, ttl time.Duration) (err error) {
	url := s.addr + setEndpoint

	body := bytes.NewReader(value)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}

	putKey(req, key)
	putTTL(req, ttl)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
func (s *Shard) Delete(ctx context.Context, key string) (err error) {
	url := s.addr + deleteEndpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, http.NoBody)
	if err != nil {
		return err
	}

	putKey(req, key)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (s *Shard) IsAlive() bool {
	return s.alive
}
func (s *Shard) Addr() string {
	return s.addr
}

func (s *Shard) Run() {
	var alive bool
	s.alive = s.healthCheck()
	slog.Info(fmt.Sprintf("health check status %q is alive: %v", s.addr, s.alive))

	go func() {
		for {
			time.Sleep(s.healthCheckInterval)
			alive = s.healthCheck()
			if s.alive && !alive {
				slog.Error(fmt.Sprintf("storage %q is not alive anymore", s.addr))
			} else if !s.alive && alive {
				slog.Info(fmt.Sprintf("storage %q is alive", s.addr))
			}
			s.alive = alive
		}
	}()
}

func putKey(req *http.Request, key string) {
	query := req.URL.Query()
	query.Set("key", key)
	req.URL.RawQuery = query.Encode()
}

func putTTL(req *http.Request, ttl time.Duration) {
	query := req.URL.Query()
	query.Set("ttl", ttl.String())
	req.URL.RawQuery = query.Encode()
}

func (s *Shard) healthCheck() bool {
	url := s.addr + healthCheckEndpoint
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		slog.Error(fmt.Sprintf("health check %q failed: %v", s.addr, err))
		return false
	}

	_, err = s.client.Do(req)
	if err != nil {
		return false
	}
	return true
}
