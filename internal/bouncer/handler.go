package bouncer

import (
	"context"
	"fmt"
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
	Get(ctx context.Context, key string) (value []byte, err error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) (err error)
	Delete(ctx context.Context, key string) (err error)
}

func (h *Handler) GetHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	key := query.Get("key")
	if key == "" {
		errorHandle(ctx, w, "key query param cant be empty", http.StatusBadRequest)
		return
	}

	value, err := h.s.Get(ctx, key)
	if err != nil {
		errorHandle(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(value)
}

func (h *Handler) SetHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	key := query.Get("key")
	if key == "" {
		errorHandle(ctx, w, "key query param cant be empty", http.StatusBadRequest)
		return
	}

	ttlStr := query.Get("ttl")
	if ttlStr == "" {
		ttlStr = "0"
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		errorHandle(ctx, w, fmt.Sprintf("cant parse ttl %q to time duration: %v", ttlStr, err), http.StatusBadRequest)
		return
	}

	value, err := io.ReadAll(r.Body)
	if err != nil {
		errorHandle(ctx, w, fmt.Sprintf("cant read value from body: %v", err), http.StatusInternalServerError)
		return
	}

	err = h.s.Set(r.Context(), key, value, ttl)
	if err != nil {
		errorHandle(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	key := query.Get("key")
	if key == "" {
		errorHandle(ctx, w, "key query param cant be empty", http.StatusBadRequest)
		return
	}

	err := h.s.Delete(ctx, key)
	if err != nil {
		errorHandle(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func errorHandle(ctx context.Context, w http.ResponseWriter, errText string, code int) {
	slog.ErrorContext(ctx, errText)
	http.Error(w, errText, code)
}
