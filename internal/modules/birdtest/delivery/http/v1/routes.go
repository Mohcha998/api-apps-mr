package v1

import (
	"apps/config"
	"apps/internal/infrastructure/db"
	"apps/internal/modules/birdtest/repository"
	"apps/internal/modules/birdtest/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, cfg *config.Config, mysql db.MysqlDBInterface) {
	repo := repository.NewBirdtestRepository(mysql)
	uc := usecase.NewBirdtestUsecase(repo, cfg.App.Timeout*time.Second)
	handler := NewBirdtestHandler(uc)

	b := r.Group("/birdtest")

	// PUT /birdtest/update_status
	b.Put("/update_status", handler.UpdateStatus)

	// POST /birdtest/status
	b.Post("/status", handler.GetStatusByEmail)
}
