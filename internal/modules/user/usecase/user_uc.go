package usecase

import (
	"apps/internal/domain"
	"apps/internal/transport/req"
	"apps/utils/helper"
	"apps/utils/response"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserUsecase struct {
	userRepo      domain.UserRepository
	ctxTimeout     time.Duration
}

func NewUserUsecase(userRepo domain.UserRepository, ctxTimeout time.Duration) *UserUsecase {
	return &UserUsecase{
		userRepo:      userRepo,
		ctxTimeout:       ctxTimeout,
	}
}

func (u *UserUsecase) Create(ctx context.Context, req *req.UserCreateReq) (*domain.User, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	_, err := u.userRepo.GetByEmail(c, req.Email)
	if err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrRegisterConflict
		}
		return nil, response.ErrRegisterConflict
	}

	var user domain.User
	helper.Copy(&user, &req)

	hash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])

	if err := u.userRepo.Create(c, &user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, response.ErrRegisterConflict
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserUsecase) Login(ctx context.Context, req *req.UserLoginReq) (*domain.User, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	var user domain.User
	helper.Copy(&user, &req)

	result, err := u.userRepo.Login(c, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrNotFound
		}
		return nil, err
	}

	hash := md5.Sum([]byte(req.Password))
	hashedInputPassword := hex.EncodeToString(hash[:])

	if result.Password != hashedInputPassword {
		return nil, response.ErrUnauthorized
	}

	return result, nil
}

func (u *UserUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(c, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	user, err := u.userRepo.GetByPhone(c, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) Update(ctx context.Context, email string, req *req.UserPasswordReq) (*domain.User, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(c, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrNotFound
		}
		return nil, err
	}

	helper.Copy(&user, &req)

	hash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])
	
	if err := u.userRepo.Update(c, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, response.ErrRegisterConflict
		}
		return nil, err
	}

	return user, nil
}