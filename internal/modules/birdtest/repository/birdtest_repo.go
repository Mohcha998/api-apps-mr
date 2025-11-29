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
	result := r.db.Conn().WithContext(ctx).
		Table("user").
		Where("id_user = ?", userID).
		Update("is_birdtest", 1)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user tidak ditemukan")
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

