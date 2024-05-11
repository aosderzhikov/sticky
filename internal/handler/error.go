package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrEmptyParam   error = errors.New("key query param cannot be empty")
	ErrInvalidParam error = errors.New("invalid ttl query param")

	ErrBodyRead error = errors.New("cant read value from body")
)

func ErrorHandle(ctx context.Context, w http.ResponseWriter, err error, code int) {
	errText := err.Error()
	slog.ErrorContext(ctx, errText)
	http.Error(w, errText, code)
}
