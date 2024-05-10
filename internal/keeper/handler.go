package keeper

import (
	"io"
	"log/slog"
	"net/http"
	"time"
)

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

type Handler struct {
	s Service
}

type Service interface {
	Get(key string) (value []byte)
	Set(key string, value []byte, ttl time.Duration)
	Delete(key string)
}

func (h *Handler) GetHandle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	_, _ = w.Write(h.s.Get(key))
}

func (h *Handler) SetHandle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")

	ttlStr := query.Get("ttl")
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	value, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
	}

	h.s.Set(key, value, ttl)
}

func (h *Handler) DeletHandle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	h.s.Delete(key)
}

func (h *Handler) HealthCheckHandle(w http.ResponseWriter, r *http.Request) {}
