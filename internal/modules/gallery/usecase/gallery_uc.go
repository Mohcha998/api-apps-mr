package usecase

import (
	"apps/internal/domain"
	"apps/internal/modules/gallery/repository"
	"context"
)

type GalleryUsecase interface {
	GetGrouped(ctx context.Context) (map[string][]domain.QuoteGallery, error)
}

type galleryUC struct {
	repo repository.GalleryRepository
}

func NewGalleryUsecase(repo repository.GalleryRepository) GalleryUsecase {
	return &galleryUC{repo: repo}
}

func (u *galleryUC) GetGrouped(ctx context.Context) (map[string][]domain.QuoteGallery, error) {

	data, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := map[string][]domain.QuoteGallery{
		"Love":     {},
		"Grateful": {},
		"Family":   {},
		"Faith":    {},
		"Dream":    {},
	}

	for _, item := range data {
		if _, ok := result[item.Category]; ok {
			result[item.Category] = append(result[item.Category], item)
		}
	}

	return result, nil
}
