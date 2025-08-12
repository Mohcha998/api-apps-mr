package domain

import (
	"context"
	"time"
)

type Youtube map[string]interface{}

type YoutubeRepository interface {
	GetActivity(ctx context.Context) (*Youtube, error)
	GetLatest(ctx context.Context) (*Youtube, error)
	GetRecent(ctx context.Context) (*Youtube, error)
	GetPlaylists(ctx context.Context) (*Youtube, error)
	GetPlaylistItems(ctx context.Context, playlistId string) (*Youtube, error)
	FetchOrCache(ctx context.Context, cacheKey string, reqURL string, ttl time.Duration) (*Youtube, error)
	BuildURL(path string, params map[string]string) (string)
}

type YouTubeUsecase interface {
	GetActivity(ctx context.Context) (*Youtube, error)
	GetLatest(ctx context.Context) (*Youtube, error)
	GetRecent(ctx context.Context) (*Youtube, error)
	GetPlaylists(ctx context.Context) (*Youtube, error)
	GetPlaylistItems(ctx context.Context, playlistId string) (*Youtube, error)
}