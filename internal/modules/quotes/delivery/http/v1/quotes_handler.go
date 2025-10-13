package v1

import (
	"apps/internal/domain"
	"apps/internal/modules/quotes/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type QuotesHandler struct {
	uc usecase.QuotesUsecase
}

// Constructor
func NewQuotesHandler(uc usecase.QuotesUsecase) *QuotesHandler {
	return &QuotesHandler{uc: uc}
}

// ===========================
// 📌 List Quotes
// ===========================
func (h *QuotesHandler) ListQuotes(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	pools, err := h.uc.ListPools(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": pools,
	})
}

// ===========================
// 📌 Get Today's Quote
// ===========================
func (h *QuotesHandler) GetTodayQuote(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user_id",
		})
	}

	pool, err := h.uc.GetOrAssignToday(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": pool,
	})
}

// ===========================
// 📌 Create Quote
// ===========================
func (h *QuotesHandler) CreateQuote(c *fiber.Ctx) error {
	var req domain.QuotePool
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.uc.CreatePool(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "quote created successfully",
		"data":    req,
	})
}

// ===========================
// 📌 Update Quote
// ===========================
func (h *QuotesHandler) UpdateQuote(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid quote ID",
		})
	}

	var req domain.QuotePool
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	req.ID = id

	changedBy := int64(1) // contoh: user ID admin, nanti bisa ambil dari JWT

	if err := h.uc.UpdatePool(c.Context(), &req, changedBy); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "quote updated successfully",
	})
}
