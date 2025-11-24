package v1

import (
	adsUsecase "apps/internal/modules/ads/usecase"

	"github.com/gofiber/fiber/v2"
)

type AdsHandler struct {
	uc adsUsecase.AdsUsecase
}

func NewAdsHandler(uc adsUsecase.AdsUsecase) *AdsHandler {
	return &AdsHandler{uc: uc}
}

func (h *AdsHandler) GetPopup(c *fiber.Ctx) error {
	result, err := h.uc.GetPopup(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "ads not found",
		})
	}

	return c.JSON(result)
}
