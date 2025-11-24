package v1

import (
    "apps/internal/modules/article/usecase"
    "context"

    "github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
    uc usecase.ArticleUsecase
}

func NewArticleHandler(uc usecase.ArticleUsecase) *ArticleHandler {
    return &ArticleHandler{uc: uc}
}

func (h *ArticleHandler) GetLatest(c *fiber.Ctx) error {
    data, err := h.uc.GetLatest(context.Background())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status": false,
            "error":  err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "status": true,
        "data":   data.Data,
    })
}
