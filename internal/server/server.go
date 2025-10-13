package server

import (
	"apps/config"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	userRoutes "apps/internal/modules/user/delivery/http/v1"
	youtubeRoutes "apps/internal/modules/youtube/delivery/http/v1"
	quotesRoutes "apps/internal/modules/quotes/delivery/http/v1"
	"apps/utils/response"

	"github.com/gofiber/fiber/v2"
)

type server struct {
	fiber *fiber.App
	cfg   *config.Config
	db    db.MysqlDBInterface
	redis *cache.Client // 🔹 sudah pakai *cache.Client
}

// NewHttpServer membuat instance server baru dengan dependency injection
func NewHttpServer(cfg *config.Config, db db.MysqlDBInterface, redis *cache.Client) *server {
	srv := &server{
		fiber: fiber.New(fiber.Config{
			ErrorHandler: response.ErrorHandler,
		}),
		cfg:   cfg,
		db:    db,
		redis: redis,
	}

	v1 := srv.fiber.Group("/api/v1")

	// Register YouTube routes
	youtubeRoutes.Routes(v1, cfg, db, redis)

	// Register User routes
	userRoutes.Routes(v1, cfg, db)

	// Register Quotes routes
	quotesRoutes.Routes(v1, cfg, db, redis)

	return srv
}

// Jalankan server
func (s *server) Run() error {
	s.fiber.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"app-name":    s.cfg.App.Name,
			"app-version": s.cfg.App.Version,
			"environment": s.cfg.App.Environment,
			"app-timeout": s.cfg.App.Timeout,
		})
	})

	return s.fiber.Listen(s.cfg.App.Port)
}
