package v1

import (
	"apps/internal/modules/resource/usecase"
	"context"

	"github.com/gofiber/fiber/v2"
)

type ResourceHandler struct {
	uc usecase.ResourceUsecase
}

func NewResourceHandler(uc usecase.ResourceUsecase) *ResourceHandler {
	return &ResourceHandler{uc: uc}
}

func (h *ResourceHandler) GetResources(c *fiber.Ctx) error {
	ctx := context.Background()

	data, err := h.uc.GetActive(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to fetch resource",
		})
	}

	// Output mirip PHP CI
	return c.JSON(fiber.Map{
		"freeResources": data,
	})
}
