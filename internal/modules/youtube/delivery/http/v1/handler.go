package v1

import (
	"apps/internal/modules/youtube/usecase"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type YouTubeHandler struct {
	usecase usecase.YouTubeUsecase
}

func NewYouTubeHandler(u usecase.YouTubeUsecase) *YouTubeHandler {
	return &YouTubeHandler{usecase: u}
}

func (h *YouTubeHandler) parseForceTTL(c *fiber.Ctx) (bool, time.Duration) {
	force := c.Query("force") == "true"
	ttl := 86400 * time.Second
	if t := c.Query("ttl"); t != "" {
		if i, err := strconv.Atoi(t); err == nil && i > 0 {
			ttl = time.Duration(i) * time.Second
		}
	}
	return force, ttl
}

func (h *YouTubeHandler) Activity(c *fiber.Ctx) error {
	force, ttl := h.parseForceTTL(c)
	data, err := h.usecase.GetActivity(context.Background(), force, ttl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func (h *YouTubeHandler) Latest(c *fiber.Ctx) error {
	force, ttl := h.parseForceTTL(c)
	data, err := h.usecase.GetLatest(context.Background(), force, ttl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func (h *YouTubeHandler) Recent(c *fiber.Ctx) error {
	force, ttl := h.parseForceTTL(c)
	data, err := h.usecase.GetRecent(context.Background(), force, ttl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func (h *YouTubeHandler) Playlists(c *fiber.Ctx) error {
	force, ttl := h.parseForceTTL(c)
	data, err := h.usecase.GetPlaylists(context.Background(), force, ttl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func (h *YouTubeHandler) PlaylistItems(c *fiber.Ctx) error {
	playlistId := c.Query("playlistId")
	if playlistId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "playlistId is required"})
	}
	force, ttl := h.parseForceTTL(c)
	data, err := h.usecase.GetPlaylistItems(context.Background(), playlistId, force, ttl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}
