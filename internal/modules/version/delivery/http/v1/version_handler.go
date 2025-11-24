package v1

import (
	"github.com/gofiber/fiber/v2"
	versionUsecase "apps/internal/modules/version/usecase"
	"apps/internal/domain"
)

type VersionHandler struct {
	uc versionUsecase.VersionUsecase
}

func NewVersionHandler(uc versionUsecase.VersionUsecase) *VersionHandler {
	return &VersionHandler{uc: uc}
}

func (h *VersionHandler) GetLatest(c *fiber.Ctx) error {
	result, err := h.uc.GetLatest(c.Context())
	if err != nil {
		return c.JSON(fiber.Map{
			"latestVersion": []interface{}{},
		})
	}

	// CI style → always array
	response := []domain.Version{}
	if result != nil {
		response = append(response, *result)
	}

	return c.JSON(fiber.Map{
		"latestVersion": response,
	})
}

