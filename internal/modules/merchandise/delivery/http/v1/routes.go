package v1

import (
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	"apps/internal/modules/merchandise/repository"
	"apps/internal/modules/merchandise/usecase"

	"github.com/gofiber/fiber/v2"
)

func Routes(
	r fiber.Router,
	db db.MysqlDBInterface,
	cache *cache.Client,
) {
	repo := repository.NewMerchandiseRepository(db, cache)
	uc := usecase.NewMerchandiseUsecase(repo)
	handler := NewMerchandiseHandler(uc)

	group := r.Group("/merchandise")
	group.Get("/", handler.GetAll)
	group.Get("/mz", handler.MZ)
	group.Get("/mrs", handler.Mrs)
	group.Get("/primerry", handler.Primerry)
	group.Get("/bytipe", handler.FindByTipe)
	group.Get("/byid", handler.FindByID)
	group.Get("/all", handler.AllMerchandise)
}
