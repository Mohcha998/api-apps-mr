package domain

import (
	"context"
	"time"
)

// Struct utama untuk response YouTube API
type Youtube struct {
	Kind  string      `json:"kind"`
	Items interface{} `json:"items"`
}

// Interface untuk layer repository
type YoutubeRepository interface {
	BuildURL(path string, params map[string]string) string
	FetchOrCache(ctx context.Context, cacheKey string, reqURL string, ttl time.Duration) (*Youtube, error)
}

// Interface untuk layer usecase
type YouTubeUsecase interface {
	GetActivity(ctx context.Context) (*Youtube, error)
	GetLatest(ctx context.Context) (*Youtube, error)
	GetRecent(ctx context.Context) (*Youtube, error)
	GetPlaylists(ctx context.Context) (*Youtube, error)
	GetPlaylistItems(ctx context.Context, playlistId string) (*Youtube, error)
}
