package usecase

import (
	"apps/internal/modules/youtube/entity"
	"apps/internal/modules/youtube/repository"
	"context"
	"time"
)

type YouTubeUsecase interface {
	GetActivity(ctx context.Context, force bool, ttl time.Duration) (entity.JSONMap, error)
	GetLatest(ctx context.Context, force bool, ttl time.Duration) (entity.JSONMap, error)
	GetRecent(ctx context.Context, force bool, ttl time.Duration) (entity.JSONMap, error)
	GetPlaylists(ctx context.Context, force bool, ttl time.Duration) (entity.JSONMap, error)
	GetPlaylistItems(ctx context.Context, playlistId string, force bool, ttl time.Duration) (entity.JSONMap, error)
}

type youtubeUsecase struct {
	repo      repository.YouTubeRepository
	httpTTL   time.Duration
}

func NewYouTubeUsecase(repo repository.YouTubeRepository, ttl time.Duration) YouTubeUsecase {
	return &youtubeUsecase{
		repo:    repo,
		httpTTL: ttl,
	}
}

func (u *youtubeUsecase) GetActivity(ctx context.Context, force bool, ttl time.Duration) (entity.JSONMap, error) {
	if ttl <= 0 {
		ttl = u.httpTTL
	}
	url := u.repo.BuildURL("activities", map[string]string{
		"part":       "snippet,contentDetails",
		"maxResults": "1",
	})
	return u.repo.FetchOrCache(ctx, "youtube_activity", url, ttl, force)
}

func (u *youtubeUsecase) GetLatest(ctx context.Context, force bool, ttl time.Duration) (entity.JSONMap, error) {
	if ttl <= 0 {
		ttl = u.httpTTL
	}
	url := u.repo.BuildURL("search", map[string]string{
		"part":       "snippet",
		"order":      "date",
		"type":       "video",
		"maxResults": "5",
	})
	return u.repo.FetchOrCache(ctx, "youtube_latest", url, ttl, force)
}

func (u *youtubeUsecase) GetRecent(ctx context.Context, force bool, ttl time.Duration) (entity.JSONMap, error) {
	if ttl <= 0 {
		ttl = u.httpTTL
	}
	url := u.repo.BuildURL("activities", map[string]string{
		"part":       "snippet,contentDetails",
		"maxResults": "5",
	})
	return u.repo.FetchOrCache(ctx, "youtube_recent", url, ttl, force)
}

func (u *youtubeUsecase) GetPlaylists(ctx context.Context, force bool, ttl time.Duration) (entity.JSONMap, error) {
	if ttl <= 0 {
		ttl = u.httpTTL
	}
	url := u.repo.BuildURL("playlists", map[string]string{
		"part":       "snippet,id,status,contentDetails",
		"maxResults": "20",
	})
	return u.repo.FetchOrCache(ctx, "youtube_playlists", url, ttl, force)
}

func (u *youtubeUsecase) GetPlaylistItems(ctx context.Context, playlistId string, force bool, ttl time.Duration) (entity.JSONMap, error) {
	if ttl <= 0 {
		ttl = u.httpTTL
	}
	cacheKey := "youtube_playlist_items_" + playlistId
	url := u.repo.BuildURL("playlistItems", map[string]string{
		"part":       "snippet,status,id",
		"maxResults": "50",
		"playlistId": playlistId,
	})
	return u.repo.FetchOrCache(ctx, cacheKey, url, ttl, force)
}
