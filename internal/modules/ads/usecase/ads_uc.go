package usecase

import (
	"apps/internal/domain"
	adsRepo "apps/internal/modules/ads/repository"
	"context"
	"time"
)

type AdsUsecase interface {
	GetPopup(ctx context.Context) (*domain.AdsPopup, error)
}

type adsUsecase struct {
	repo    adsRepo.AdsRepository
	timeout time.Duration
}

func NewAdsUsecase(repo adsRepo.AdsRepository, timeout time.Duration) AdsUsecase {
	return &adsUsecase{
		repo:    repo,
		timeout: timeout,
	}
}

func (u *adsUsecase) GetPopup(ctx context.Context) (*domain.AdsPopup, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.GetActivePopup(ctx)
}
