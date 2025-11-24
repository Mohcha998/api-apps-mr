package v1

import (
	"apps/config"
	"apps/internal/infrastructure/db"

	versionRepo "apps/internal/modules/version/repository"
	versionUsecase "apps/internal/modules/version/usecase"

	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, cfg *config.Config, db db.MysqlDBInterface) {
	repo := versionRepo.NewVersionRepository(db)
	uc := versionUsecase.NewVersionUsecase(repo)
	handler := NewVersionHandler(uc)

	r.Get("/version", handler.GetLatest)
}
