package handler

import (
	"errors"
	"net/http"
	"time"
)

const (
	keyParam = "key"
	ttlParam = "ttl"
)

func ExtractKeyAndTTL(r *http.Request) (key string, ttl time.Duration, err error) {
	key, err = ExtractKey(r)
	if err != nil {
		return "", 0, err
	}

	query := r.URL.Query()
	ttlStr := query.Get("ttl")
	if ttlStr == "" {
		ttlStr = "0"
	}
	ttl, err = time.ParseDuration(ttlStr)
	if err != nil {
		return "", 0, errors.Join(ErrInvalidParam, err)
	}

	return key, ttl, nil
}

func ExtractKey(r *http.Request) (string, error) {
	query := r.URL.Query()
	key := query.Get(keyParam)
	if key == "" {
		return "", ErrEmptyParam
	}
	return key, nil
}
