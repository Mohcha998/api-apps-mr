package v1

import (
	"apps/internal/domain"
	"apps/internal/transport/req"
	"apps/utils/response"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *userHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) Create(c *fiber.Ctx) error {
	var request req.UserCreateReq
	var ctx = c.Context()

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	user, err := h.userUsecase.Create(ctx, &request)
	if err != nil {
		return err
	}

	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, user))
}

func (h *userHandler) Update(c *fiber.Ctx) error {
	var request req.UserPasswordReq
	var ctx = c.Context()
	var email = c.Params("email")

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	product, err := h.userUsecase.Update(ctx, email, &request)
	if err != nil {
		return err
	}

	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, product))
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	var request req.UserLoginReq
	var ctx = c.Context()

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	user, err := h.userUsecase.Login(ctx, &request)
	if err != nil {
		return err
	}

	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, user))
}

func (h *userHandler) GetByEmail(c *fiber.Ctx) error {
	var ctx = c.Context()
	var email = c.Params("email")

	user, err := h.userUsecase.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, user))
}

func (h *userHandler) GetByPhone(c *fiber.Ctx) error {
	var ctx = c.Context()
	var phone = c.Params("phone")

	user, err := h.userUsecase.GetByPhone(ctx, phone)
	if err != nil {
		return err
	}

	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, user))
}