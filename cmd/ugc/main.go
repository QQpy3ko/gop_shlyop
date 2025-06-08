package main

import (
	"context"
	"gop_shlyop/internal/config"
	"gop_shlyop/internal/logger"
	"gop_shlyop/internal/server"
	"os"
)

func main() {
	// load config
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	cfg := config.MustLoad(configPath)

	log := logger.New(cfg.Env)

	log.Info().Str("env", cfg.Env).Msg("UGC Service is starting")
	log.Debug().Msg("debug messages are enabled")

	ctx := context.Background()

    // create server
	srv, err := server.New(ctx, cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
	}

	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run server")
	}
}
