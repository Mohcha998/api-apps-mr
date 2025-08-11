package v1

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, handler *YouTubeHandler) {
	youtube := r.Group("/youtube")

	youtube.Get("/activity", handler.Activity)
	youtube.Get("/latest", handler.Latest)
	youtube.Get("/recent", handler.Recent)
	youtube.Get("/playlists", handler.Playlists)
	youtube.Get("/playlist-items", handler.PlaylistItems)
}
