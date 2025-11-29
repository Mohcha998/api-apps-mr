package v1

import (
	"apps/internal/modules/birdtest/usecase"
	"github.com/gofiber/fiber/v2"
	// "strconv"
)

type BirdtestHandler struct {
	uc usecase.BirdtestUsecase
}

func NewBirdtestHandler(uc usecase.BirdtestUsecase) *BirdtestHandler {
	return &BirdtestHandler{uc: uc}
}

func (h *BirdtestHandler) UpdateStatus(ctx *fiber.Ctx) error {
	var body struct {
		UserID int `json:"user_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid JSON",
		})
	}

	if body.UserID == 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "User ID harus diisi",
		})
	}

	if err := h.uc.UpdateStatus(ctx.Context(), body.UserID); err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  true,
		"message": "Status Birdtest berhasil diperbarui",
	})
}

func (h *BirdtestHandler) GetStatusByEmail(ctx *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid JSON",
		})
	}

	res, err := h.uc.GetStatusByEmail(ctx.Context(), body.Email)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Email tidak ditemukan",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":      true,
		"is_birdtest": res.IsBirdtest,
	})
}
