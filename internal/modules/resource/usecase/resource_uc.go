package usecase

import (
    "apps/internal/domain"
    "apps/internal/modules/resource/repository"
    "context"
)

type ResourceUsecase interface {
    GetActive(ctx context.Context) ([]domain.Resource, error)
}

type resourceUC struct {
    repo repository.ResourceRepository
}

func NewResourceUsecase(repo repository.ResourceRepository) ResourceUsecase {
    return &resourceUC{repo: repo}
}

func (uc *resourceUC) GetActive(ctx context.Context) ([]domain.Resource, error) {
    return uc.repo.GetActive(ctx)
}
