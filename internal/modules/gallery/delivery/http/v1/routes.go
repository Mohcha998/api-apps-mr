package v1

import (
	"apps/internal/infrastructure/db"
	"apps/internal/modules/gallery/repository"
	"apps/internal/modules/gallery/usecase"

	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, db db.MysqlDBInterface) {
	repo := repository.NewGalleryRepository(db)
	uc := usecase.NewGalleryUsecase(repo)
	handler := NewGalleryHandler(uc)

	group := r.Group("/gallery")
	group.Get("/", handler.GetAllGrouped)
}
