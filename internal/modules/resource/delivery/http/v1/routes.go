package v1

import (
	"apps/config"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	"apps/internal/modules/resource/repository"
	"apps/internal/modules/resource/usecase"

	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, cfg *config.Config, db db.MysqlDBInterface, redis *cache.Client) {
    repo := repository.NewResourceRepository(db, redis)
    uc := usecase.NewResourceUsecase(repo)
    handler := NewResourceHandler(uc)

    group := r.Group("/resources")
    group.Get("/", handler.GetResources)
}

