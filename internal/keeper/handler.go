//go:generate mockgen -source=$GOFILE -destination=mock_test.go -package=$GOPACKAGE
package keeper

import (
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
	Get(key string) (value []byte)
	Set(key string, value []byte, ttl time.Duration)
	Delete(key string)
}

func (h *Handler) GetHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key, err := handler.ExtractKey(r)
	if err != nil {
		handler.ErrorHandle(ctx, w, err, http.StatusBadRequest)
		return
	}

	value := h.s.Get(key)
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

	h.s.Set(key, value, ttl)
}

func (h *Handler) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key, err := handler.ExtractKey(r)
	if err != nil {
		handler.ErrorHandle(ctx, w, err, http.StatusBadRequest)
		return
	}

	h.s.Delete(key)
}

func (h *Handler) HealthCheckHandle(w http.ResponseWriter, r *http.Request) {}
