package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aosderzhikov/sticky/internal/keeper"
)

const (
	httpAddrEnv = "HTTP_ADDRESS"
	defaultAddr = "localhost:8181"
)

func main() {
	addr := os.Getenv(httpAddrEnv)
	if addr == "" {
		addr = defaultAddr
	}

	k := keeper.NewService()
	k.Run()

	handler := keeper.NewHandler(k)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /get", handler.GetHandle)
	mux.HandleFunc("POST /set", handler.SetHandle)
	mux.HandleFunc("DELETE /delete", handler.DeletHandle)

	slog.Info(fmt.Sprintf("start keeper on %q", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error(err.Error())
	}
}
