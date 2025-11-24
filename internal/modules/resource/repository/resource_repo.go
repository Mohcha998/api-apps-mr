package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/cache"
	"apps/internal/infrastructure/db"
	"context"
	"time"
)

type ResourceRepository interface {
	GetActive(ctx context.Context) ([]domain.Resource, error)
}

type resourceRepo struct {
	db    db.MysqlDBInterface
	cache *cache.Client
}

func NewResourceRepository(db db.MysqlDBInterface, cache *cache.Client) ResourceRepository {
	return &resourceRepo{db: db, cache: cache}
}

func (r *resourceRepo) GetActive(ctx context.Context) ([]domain.Resource, error) {
    key := "resource:list"

    var cached []domain.Resource
    if err := r.cache.GetJSON(ctx, key, &cached); err == nil && len(cached) > 0 {
        return cached, nil
    }

    var result []domain.Resource

    query := db.NewQuery("status = ?", "1")

    if err := r.db.Find(ctx, &result, db.WithQuery(query)); err != nil {
        return nil, err
    }

    _ = r.cache.SetJSON(ctx, key, result, 10*time.Minute)

    return result, nil
}

