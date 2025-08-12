package usecase

import (
	"apps/internal/domain"
	"context"
	"time"
)



type youtubeUsecase struct {
	youtubeRepo  domain.YoutubeRepository
	ctxTimeout   time.Duration
}

func NewYouTubeUsecase(youtubeRepo domain.YoutubeRepository, ctxTimeout time.Duration) *youtubeUsecase {
	return &youtubeUsecase{
		youtubeRepo: youtubeRepo,
		ctxTimeout: ctxTimeout,
	}
}

func (u *youtubeUsecase) GetActivity(ctx context.Context) (*domain.Youtube, error) {
	url := u.youtubeRepo.BuildURL("activities", map[string]string{
		"part":       "snippet,contentDetails",
		"maxResults": "1",
	})
	return u.youtubeRepo.FetchOrCache(ctx, "youtube_activity", url, 60*time.Second)
}

func (u *youtubeUsecase) GetLatest(ctx context.Context) (*domain.Youtube, error) {

	url := u.youtubeRepo.BuildURL("search", map[string]string{
		"part":       "snippet",
		"order":      "date",
		"type":       "video",
		"maxResults": "5",
	})
	return u.youtubeRepo.FetchOrCache(ctx, "youtube_latest", url, 60*time.Second)
}

func (u *youtubeUsecase) GetRecent(ctx context.Context) (*domain.Youtube, error) {
	
	url := u.youtubeRepo.BuildURL("activities", map[string]string{
		"part":       "snippet,contentDetails",
		"maxResults": "5",
	})
	return u.youtubeRepo.FetchOrCache(ctx, "youtube_recent", url, 60*time.Second)
}

func (u *youtubeUsecase) GetPlaylists(ctx context.Context) (*domain.Youtube, error) {
	url := u.youtubeRepo.BuildURL("playlists", map[string]string{
		"part":       "snippet,id,status,contentDetails",
		"maxResults": "20",
	})
	return u.youtubeRepo.FetchOrCache(ctx, "youtube_playlists", url, 60*time.Second)
}

func (u *youtubeUsecase) GetPlaylistItems(ctx context.Context, playlistId string) (*domain.Youtube, error) {
	cacheKey := "youtube_playlist_items_" + playlistId
	url := u.youtubeRepo.BuildURL("playlistItems", map[string]string{
		"part":       "snippet,status,id",
		"maxResults": "50",
		"playlistId": playlistId,
	})
	return u.youtubeRepo.FetchOrCache(ctx, cacheKey, url, 60*time.Second)
}
