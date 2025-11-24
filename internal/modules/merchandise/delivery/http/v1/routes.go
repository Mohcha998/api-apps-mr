package v1

import (
	"apps/internal/modules/merchandise/usecase"
	"apps/internal/modules/merchandise/repository"
	"apps/internal/infrastructure/db"

	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, db db.MysqlDBInterface) {
	repo := repository.NewMerchandiseRepository(db)
	uc := usecase.NewMerchandiseUsecase(repo)
	handler := NewMerchandiseHandler(uc)

	group := r.Group("/merchandise")
	group.Get("/", handler.GetAll)
	group.Get("/mz", handler.MZ)
	group.Get("/primerry", handler.Primerry)
	group.Get("/bytipe", handler.FindByTipe)
	group.Get("/byid", handler.FindByID)
}
