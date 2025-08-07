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
		log.Fatal(err)
	}

	mysql, err := db.NewMysqlConnection(cfg)
	if err != nil {
		log.Println(err)
	}

	redis := cache.NewRedisConnection(cfg)

	httpServer := server.NewHttpServer(cfg, mysql, redis)
	if err = httpServer.Run(); err != nil {
		log.Fatal(err)
	}
}
