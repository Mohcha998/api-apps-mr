package v1

import (
    "apps/config"
    "apps/internal/infrastructure/cache"
    "apps/internal/infrastructure/db"

    articleRepo "apps/internal/modules/article/repository"
    articleUsecase "apps/internal/modules/article/usecase"
    "github.com/gofiber/fiber/v2"
)

func Routes(router fiber.Router, cfg *config.Config, db db.MysqlDBInterface, redis *cache.Client) {
    repo := articleRepo.NewArticleRepository(redis)

    uc := articleUsecase.NewArticleUsecase(repo)

    handler := NewArticleHandler(uc)

    r := router.Group("/articles")
    r.Get("/latest", handler.GetLatest)
}
