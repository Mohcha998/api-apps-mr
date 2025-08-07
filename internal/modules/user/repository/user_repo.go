package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/db"
	"context"
	"errors"
)

type userRepository struct {
	db db.MysqlDBInterface
}

func NewUserRepository(db db.MysqlDBInterface) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Login(ctx context.Context, user *domain.User) (*domain.User, error) {
	var result domain.User
	var query db.Query
	if user.Email != "" {
		query = db.NewQuery("email = ?", user.Email)
	} else if user.Mobile != "" {
		query = db.NewQuery("mobile = ?", user.Mobile)
	} else {
		return nil, errors.New("email or phone must be provided")
	}
	
	if err := r.db.FindOne(ctx, &result, db.WithQuery(query), db.WithOrder("date_created DESC")); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.Create(ctx, user)
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.Update(ctx, user)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := db.NewQuery("email = ?", email)
	if err := r.db.FindOne(ctx, &user, db.WithQuery(query), db.WithOrder("date_created DESC")); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, mobile string) (*domain.User, error) {
	var user domain.User
	query := db.NewQuery("mobile = ?", mobile)
	if err := r.db.FindOne(ctx, &user, db.WithQuery(query), db.WithOrder("date_created DESC")); err != nil {
		return nil, err
	}

	return &user, nil
}

