package v1

import (
	"apps/internal/modules/gallery/usecase"

	"github.com/gofiber/fiber/v2"
)

type GalleryHandler struct {
	uc usecase.GalleryUsecase
}

func NewGalleryHandler(uc usecase.GalleryUsecase) *GalleryHandler {
	return &GalleryHandler{uc: uc}
}

func (h *GalleryHandler) GetAllGrouped(c *fiber.Ctx) error {
	res, err := h.uc.GetGrouped(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": true,
		"data":   res,
	})
}
