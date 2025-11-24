package v1

import (
	"apps/config"
	"apps/internal/infrastructure/db"

	adsRepo "apps/internal/modules/ads/repository"
	adsUsecase "apps/internal/modules/ads/usecase"

	"time"

	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, cfg *config.Config, db db.MysqlDBInterface) {
	repo := adsRepo.NewAdsRepository(db)
	uc := adsUsecase.NewAdsUsecase(repo, cfg.App.Timeout*time.Second)
	handler := NewAdsHandler(uc)

	r.Get("/popup", handler.GetPopup)
}
