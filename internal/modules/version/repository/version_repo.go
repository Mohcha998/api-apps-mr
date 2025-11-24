package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/db"
	"context"
	"database/sql"
)

type VersionRepository interface {
	GetLatest(ctx context.Context) (*domain.Version, error)
}

type versionRepository struct {
	db db.MysqlDBInterface
}

func NewVersionRepository(db db.MysqlDBInterface) VersionRepository {
	return &versionRepository{db: db}
}

func (r *versionRepository) GetLatest(ctx context.Context) (*domain.Version, error) {
	var version domain.Version

	err := r.db.FindOne(
		ctx,
		&version,
		db.WithOrder("id DESC"),
		db.WithLimit(1),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &version, nil
}
