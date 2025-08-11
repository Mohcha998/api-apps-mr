package main

import (
	"apps/config"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	"apps/internal/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mysql, err := db.NewMysqlConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect mysql: %v", err)
	}

	redis := cache.NewRedisConnection(cfg)
	if redis == nil {
		log.Fatalf("failed to connect redis")
	}

	httpServer := server.NewHttpServer(cfg, mysql, redis)
	if err = httpServer.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
