package v1

import (
	"apps/config"
	"apps/internal/infrastructure/db"
	youtubeRepo "apps/internal/modules/youtube/repository"
	youtubeUsecase "apps/internal/modules/youtube/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func Routes(r fiber.Router, cfg *config.Config, db db.MysqlDBInterface, redis *redis.Client) {
	
	youtubeRepo := youtubeRepo.NewYouTubeRepository(cfg.YouTube.APIKey, cfg.YouTube.ChannelID, time.Second*10, redis)
	youtubeUC := youtubeUsecase.NewYouTubeUsecase(youtubeRepo, cfg.App.Timeout*time.Second)
	youtubeHandler := NewYouTubeHandler(youtubeUC)

	youtube := r.Group("/youtube")

	youtube.Get("/activity", youtubeHandler.Activity)
	youtube.Get("/latest", youtubeHandler.Latest)
	youtube.Get("/recent", youtubeHandler.Recent)
	youtube.Get("/playlists", youtubeHandler.Playlists)
	youtube.Get("/playlist-items", youtubeHandler.PlaylistItems)
}
