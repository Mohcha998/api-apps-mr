package usecase

import (
    "context"
    "apps/internal/domain"
)

type ArticleUsecase interface {
    GetLatest(ctx context.Context) (domain.LatestArticles, error)
}

type articleUsecase struct {
    repo domain.ArticleRepository
}

func NewArticleUsecase(r domain.ArticleRepository) ArticleUsecase {
    return &articleUsecase{repo: r}
}

func (uc *articleUsecase) GetLatest(ctx context.Context) (domain.LatestArticles, error) {
    return uc.repo.GetLatest(ctx)
}
