package usecase

import (
	"apps/internal/domain"
	"context"
	"time"
)

type youtubeUsecase struct {
	youtubeRepo domain.YoutubeRepository
	ctxTimeout  time.Duration
}

func NewYouTubeUsecase(youtubeRepo domain.YoutubeRepository, ctxTimeout time.Duration) domain.YouTubeUsecase {
	return &youtubeUsecase{
		youtubeRepo: youtubeRepo,
		ctxTimeout:  ctxTimeout,
	}
}

func (u *youtubeUsecase) GetActivity(ctx context.Context) (*domain.Youtube, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()
	return u.youtubeRepo.GetActivity(ctx)
}

func (u *youtubeUsecase) GetLatest(ctx context.Context) (*domain.Youtube, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()
	return u.youtubeRepo.GetLatest(ctx)
}

func (u *youtubeUsecase) GetRecent(ctx context.Context) (*domain.Youtube, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()
	return u.youtubeRepo.GetRecent(ctx)
}

func (u *youtubeUsecase) GetPlaylists(ctx context.Context) (*domain.Youtube, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()
	return u.youtubeRepo.GetPlaylists(ctx)
}

func (u *youtubeUsecase) GetPlaylistItems(ctx context.Context, playlistId string) (*domain.Youtube, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()
	return u.youtubeRepo.GetPlaylistItems(ctx, playlistId)
}
