package usecase

import (
	"apps/internal/domain"
	"apps/internal/modules/birdtest/repository"
	"context"
	"errors"
	"time"
)

type BirdtestUsecase interface {
	UpdateStatus(ctx context.Context, userID int) error
	GetStatusByEmail(ctx context.Context, email string) (*domain.BirdtestStatus, error)
}

type birdtestUsecase struct {
	repo    repository.BirdtestRepository
	timeout time.Duration
}

func NewBirdtestUsecase(repo repository.BirdtestRepository, timeout time.Duration) BirdtestUsecase {
	return &birdtestUsecase{
		repo:    repo,
		timeout: timeout,
	}
}

func (u *birdtestUsecase) UpdateStatus(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.UpdateStatus(ctx, userID)
}

func (u *birdtestUsecase) GetStatusByEmail(ctx context.Context, email string) (*domain.BirdtestStatus, error) {
	if email == "" {
		return nil, errors.New("email harus diisi")
	}

	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.GetStatusByEmail(ctx, email)
}
