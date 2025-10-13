package v1

import (
	"apps/config"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	"apps/internal/modules/quotes/repository"
	"apps/internal/modules/quotes/usecase"

	"github.com/gofiber/fiber/v2"
)

func Routes(router fiber.Router, cfg *config.Config, db db.MysqlDBInterface, redis *cache.Client) {
	repo := repository.NewQuotesRepository(db, redis)
	uc := usecase.NewQuotesUsecase(repo)
	handler := NewQuotesHandler(uc)

	r := router.Group("/quotes")
	r.Get("/", handler.ListQuotes)
	r.Get("/today/:user_id", handler.GetTodayQuote)
	r.Post("/", handler.CreateQuote)
	r.Put("/:id", handler.UpdateQuote)
}
