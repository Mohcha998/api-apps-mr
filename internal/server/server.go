package server

import (
    "apps/config"
    "apps/internal/infrastructure/cache"
    "apps/internal/infrastructure/db"

    userRoutes "apps/internal/modules/user/delivery/http/v1"
    youtubeRoutes "apps/internal/modules/youtube/delivery/http/v1"
    quotesRoutes "apps/internal/modules/quotes/delivery/http/v1"
    articleRoutes "apps/internal/modules/article/delivery/http/v1"
    resourceRoutes "apps/internal/modules/resource/delivery/http/v1"
    adsRoutes "apps/internal/modules/ads/delivery/http/v1"
    versionRoutes "apps/internal/modules/version/delivery/http/v1"
    merchRoutes "apps/internal/modules/merchandise/delivery/http/v1"

    "apps/utils/response"

    "github.com/gofiber/fiber/v2"
)

type server struct {
    fiber *fiber.App
    cfg   *config.Config
    db    db.MysqlDBInterface
    redis *cache.Client
}

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

    youtubeRoutes.Routes(v1, cfg, db, redis)
    userRoutes.Routes(v1, cfg, db)
    quotesRoutes.Routes(v1, cfg, db, redis)
    articleRoutes.Routes(v1, cfg, db, redis)
    resourceRoutes.Routes(v1, cfg, db, redis)
    adsRoutes.Routes(v1, cfg, db)
    versionRoutes.Routes(v1, cfg, db)
    merchRoutes.Routes(v1, db)

    return srv
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

    return s.fiber.Listen(s.cfg.App.Port)
}
