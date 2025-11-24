package v1

import (
	"apps/internal/modules/merchandise/usecase"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type MerchandiseHandler struct {
	uc usecase.MerchandiseUsecase
}

func NewMerchandiseHandler(uc usecase.MerchandiseUsecase) *MerchandiseHandler {
	return &MerchandiseHandler{uc: uc}
}

func (h *MerchandiseHandler) GetAll(c *fiber.Ctx) error {
	res, err := h.uc.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(res)
}

func (h *MerchandiseHandler) MZ(c *fiber.Ctx) error {
	res, err := h.uc.GetMZ(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"merchandise_all": res,
	})
}

func (h *MerchandiseHandler) Primerry(c *fiber.Ctx) error {
	res, err := h.uc.GetPrimerry(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"merchandise_all": res,
	})
}

func (h *MerchandiseHandler) FindByTipe(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)

	res, err := h.uc.FindByTipe(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"merchandise": res,
	})
}

func (h *MerchandiseHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)

	res, err := h.uc.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"merchandise": res,
	})
}
