package v1

import (
	"apps/config"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	youtubeRepo "apps/internal/modules/youtube/repository"
	youtubeUsecase "apps/internal/modules/youtube/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Routes mendefinisikan semua endpoint YouTube API v1
func Routes(r fiber.Router, cfg *config.Config, db db.MysqlDBInterface, redis *cache.Client) {
	repo := youtubeRepo.NewYouTubeRepository(cfg.YouTube.APIKey, cfg.YouTube.ChannelID, 10*time.Second, redis)
	usecase := youtubeUsecase.NewYouTubeUsecase(repo, cfg.App.Timeout*time.Second)
	handler := NewYouTubeHandler(usecase)

	youtube := r.Group("/youtube")

	youtube.Get("/activity", handler.Activity)
	youtube.Get("/latest", handler.Latest)
	youtube.Get("/recent", handler.Recent)
	youtube.Get("/playlists", handler.Playlists)
	youtube.Get("/playlist-items", handler.PlaylistItems)
}
