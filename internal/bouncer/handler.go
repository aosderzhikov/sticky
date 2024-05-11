//go:generate mockgen -source=$GOFILE -destination=mock_test.go -package=$GOPACKAGE
package bouncer

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/aosderzhikov/sticky/internal/handler"
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

	key, err := handler.ExtractKey(r)
	if err != nil {
		handler.ErrorHandle(ctx, w, err, http.StatusBadRequest)
		return
	}

	value, err := h.s.Get(ctx, key)
	if err != nil {
		handler.ErrorHandle(ctx, w, err, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(value)
}

func (h *Handler) SetHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key, ttl, err := handler.ExtractKeyAndTTL(r)
	if err != nil {
		handler.ErrorHandle(ctx, w, err, http.StatusBadRequest)
		return
	}

	value, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.Join(handler.ErrBodyRead, err)
		handler.ErrorHandle(ctx, w, err, http.StatusInternalServerError)
		return
	}

	err = h.s.Set(ctx, key, value, ttl)
	if err != nil {
		handler.ErrorHandle(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key, err := handler.ExtractKey(r)
	if err != nil {
		handler.ErrorHandle(ctx, w, err, http.StatusBadRequest)
		return
	}

	err = h.s.Delete(ctx, key)
	if err != nil {
		handler.ErrorHandle(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
