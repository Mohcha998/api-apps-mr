package v1

import (
	"apps/internal/domain"
	"apps/utils/response"

	"github.com/gofiber/fiber/v2"
)

type YouTubeHandler struct {
	YotubeUsecase domain.YouTubeUsecase
}

func NewYouTubeHandler(YotubeUsecase domain.YouTubeUsecase) *YouTubeHandler {
	return &YouTubeHandler{YotubeUsecase: YotubeUsecase}
}


func (h *YouTubeHandler) Activity(c *fiber.Ctx) error {
	var ctx = c.Context()
	data, err := h.YotubeUsecase.GetActivity(ctx)
	if err != nil {
		return err
	}
	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, data))
}

func (h *YouTubeHandler) Latest(c *fiber.Ctx) error {
	var ctx = c.Context()
	data, err := h.YotubeUsecase.GetLatest(ctx)
	if err != nil {
		return err
	}
	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, data))
}

func (h *YouTubeHandler) Recent(c *fiber.Ctx) error {
	var ctx = c.Context()
	data, err := h.YotubeUsecase.GetRecent(ctx)
	if err != nil {
		return err
	}
	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, data))
}

func (h *YouTubeHandler) Playlists(c *fiber.Ctx) error {
	var ctx = c.Context()
	data, err := h.YotubeUsecase.GetPlaylists(ctx)
	if err != nil {
		return err
	}
	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, data))
}

func (h *YouTubeHandler) PlaylistItems(c *fiber.Ctx) error {
	var ctx = c.Context()
	playlistId := c.Query("playlistId")
	if playlistId == "" {
		return c.JSON(response.ErrNotFound)
	}
	data, err := h.YotubeUsecase.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		return err
	}
	return c.JSON(response.NewSuccessResponse(fiber.StatusOK, data))
}
