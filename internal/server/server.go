package server

import (
	"apps/config"
	"apps/internal/infrastructure/db"
	userRoutes "apps/internal/modules/user/delivery/http/v1"
	"apps/utils/response"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type server struct {
	fiber *fiber.App
	cfg   *config.Config
	db    db.MysqlDBInterface
	redis *redis.Client
}

func NewHttpServer(cfg *config.Config, db db.MysqlDBInterface, redis *redis.Client) *server {
	return &server{
		fiber: fiber.New(fiber.Config{
			ErrorHandler: response.ErrorHandler,
		}),
		cfg:   cfg,
		db:    db,
		redis: redis,
	}
}

func (s *server) Run() error {
	s.fiber.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"app-name":    s.cfg.App.Name,
			"app-version": s.cfg.App.Version,
			"environment": s.cfg.App.Environment,
			"app-timeout": s.cfg.App.Timeout,
		})
	})

	// s.fiber.Get("/swagger/*", swagger.HandlerDefault)

	v1 := s.fiber.Group("api/v1")
	// categoryRoutes.Routes(v1, s.cfg, s.db, s.redis)
	userRoutes.Routes(v1, s.cfg, s.db)
	if err := s.fiber.Listen(s.cfg.App.Port); err != nil {
		return err
	}

	return nil
}
