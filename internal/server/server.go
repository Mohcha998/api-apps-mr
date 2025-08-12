package server

import (
	"apps/config"
	"apps/internal/infrastructure/db"
	userRoutes "apps/internal/modules/user/delivery/http/v1"
	youtubeRoutes "apps/internal/modules/youtube/delivery/http/v1"
	"apps/utils/response"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// server struct menyimpan dependencies server HTTP
type server struct {
	fiber *fiber.App
	cfg   *config.Config
	db    db.MysqlDBInterface
	redis *redis.Client
}

// NewHttpServer membuat instance server baru dengan dependency injection
func NewHttpServer(cfg *config.Config, db db.MysqlDBInterface, redis *redis.Client) *server {
	srv := &server{
		fiber: fiber.New(fiber.Config{
			ErrorHandler: response.ErrorHandler,
		}),
		cfg:   cfg,
		db:    db,
		redis: redis,
	}


	// Daftarkan route group api/v1
	v1 := srv.fiber.Group("/api/v1")

	// Register routes YouTube
	youtubeRoutes.Routes(v1, cfg, db, redis)

	// Register routes User
	userRoutes.Routes(v1, cfg, db)

	return srv
}

// Run menjalankan server di port yang telah dikonfigurasi
func (s *server) Run() error {
	// Endpoint root untuk info app
	s.fiber.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"app-name":    s.cfg.App.Name,
			"app-version": s.cfg.App.Version,
			"environment": s.cfg.App.Environment,
			"app-timeout": s.cfg.App.Timeout,
		})
	})

	// Jalankan server Fiber pada port yang ditentukan di config
	return s.fiber.Listen(s.cfg.App.Port)
}
