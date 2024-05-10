package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aosderzhikov/sticky/internal/bouncer"
	"github.com/caarlos0/env/v9"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Bouncer BouncerConfig `yaml:"bouncer"`
}

type BouncerConfig struct {
	Addr       string          `yaml:"addr"`
	Storages   []StorageConfig `yaml:"storages"`
	DefaultTTL time.Duration   `yaml:"defaultTtl" envDefault:"10s"`
}

type StorageConfig struct {
	Addr                string        `yaml:"addr"`
	HealthCheckInterval time.Duration `yaml:"healthCheckInterval" envDefault:"5s"`
}

func load(r io.Reader) (Config, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

const (
	configPathEnv    = "CONFIG_PATH"
	defaulConfigPath = "./cmd/bouncer/config.yaml"
)

func main() {
	configFileName := os.Getenv(configPathEnv)
	if configFileName == "" {
		configFileName = defaulConfigPath
	}

	cfgFile, err := os.Open(configFileName)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	cfg, err := load(cfgFile)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err = env.Parse(&cfg); err != nil {
		slog.Error(err.Error())
		return
	}

	storages := make([]bouncer.Storage, 0, len(cfg.Bouncer.Storages))
	for _, s := range cfg.Bouncer.Storages {
		if err = env.Parse(&s); err != nil {
			slog.Error(err.Error())
			return
		}

		shard := bouncer.NewShard(s.Addr, s.HealthCheckInterval, nil)
		shard.Run()
		storages = append(storages, shard)
	}

	service := bouncer.NewShardService(storages)
	handler := bouncer.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /get", handler.GetHandle)
	mux.HandleFunc("POST /set", handler.SetHandle)
	mux.HandleFunc("DELETE /delete", handler.DeleteHandle)

	slog.Info(fmt.Sprintf("start bouncer on %q", cfg.Bouncer.Addr))
	if err = http.ListenAndServe(cfg.Bouncer.Addr, mux); err != nil {
		slog.Error(err.Error())
	}
}
