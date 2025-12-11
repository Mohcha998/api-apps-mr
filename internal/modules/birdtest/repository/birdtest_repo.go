package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/db"
	"context"
	"errors"
)

type BirdtestRepository interface {
	UpdateStatus(ctx context.Context, userID int) error
	GetStatusByEmail(ctx context.Context, email string) (*domain.BirdtestStatus, error)
}

type birdtestRepository struct {
	db db.MysqlDBInterface
}

func NewBirdtestRepository(db db.MysqlDBInterface) BirdtestRepository {
	return &birdtestRepository{db: db}
}

func (r *birdtestRepository) UpdateStatus(ctx context.Context, userID int) error {
	// 1. Cek apakah user ada
	var count int64
	err := r.db.Conn().WithContext(ctx).
		Table("user").
		Where("id_user = ?", userID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("user tidak ditemukan")
	}

	err = r.db.Conn().WithContext(ctx).
		Table("user").
		Where("id_user = ?", userID).
		UpdateColumn("is_birdtest", 1).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *birdtestRepository) GetStatusByEmail(ctx context.Context, email string) (*domain.BirdtestStatus, error) {
	var user domain.BirdtestUser

	err := r.db.FindOne(ctx, &user,
		db.WithQuery(db.NewQuery("email = ?", email)),
		db.WithoutOrder(),
	)
	if err != nil {
		return nil, err
	}

	return &domain.BirdtestStatus{
		IsBirdtest: user.IsBirdtest,
	}, nil
}
