package main

import (
	"log"
	"os"

	"gop_shlyop/internal/config"
	"gop_shlyop/internal/server"
)

func main() {
	// load config
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	cfg := config.MustLoad(configPath)

	// create server
	srv := server.New(cfg)

	// run server
	log.Println("UGC Service is starting...")
	if err := srv.Run(); err != nil {
		log.Fatalf("failed to run server: %s", err)
	}
}
