package usecase

import (
	"context"
	"apps/internal/domain"
	"apps/internal/modules/version/repository"
)

type VersionUsecase interface {
	GetLatest(ctx context.Context) (*domain.Version, error)
}

type versionUsecase struct {
	repo repository.VersionRepository
}

func NewVersionUsecase(r repository.VersionRepository) VersionUsecase {
	return &versionUsecase{repo: r}
}

func (u *versionUsecase) GetLatest(ctx context.Context) (*domain.Version, error) {
	return u.repo.GetLatest(ctx)
}
