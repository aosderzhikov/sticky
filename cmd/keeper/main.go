package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aosderzhikov/sticky/internal/keeper"
)

const (
	httpAddrEnv = "HTTP_ADDRESS"
	ttlEnv      = "TTL"
	debugEnv    = "DEBUG"

	defaultAddr = "localhost:8181"
	defaultTTL  = "10m"
)

func main() {
	addr := os.Getenv(httpAddrEnv)
	if addr == "" {
		addr = defaultAddr
	}

	ttlStr := os.Getenv(ttlEnv)
	if ttlStr == "" {
		ttlStr = defaultTTL
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	debugMode := os.Getenv(debugEnv)
	if debugMode == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("debug level is on")
	}

	k := keeper.NewService(ttl)
	k.Run()

	handler := keeper.NewHandler(k)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /get", handler.GetHandle)
	mux.HandleFunc("POST /set", handler.SetHandle)
	mux.HandleFunc("DELETE /delete", handler.DeleteHandle)
	mux.HandleFunc("GET /health-check", handler.HealthCheckHandle)

	srv := http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 200 * time.Millisecond,
		WriteTimeout:      1 * time.Second,
		Handler:           mux,
	}

	slog.Info(fmt.Sprintf("start keeper on %q", addr))
	if err := srv.ListenAndServe(); err != nil {
		slog.Error(err.Error())
	}
}
