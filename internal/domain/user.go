package domain

import (
	"apps/internal/transport/req"
	"context"
	"time"

	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusInProgress UserStatus = "in-progress"
	UserStatusDone       UserStatus = "done"
)

type User struct {
	ID_User       uint      `gorm:"primaryKey;autoIncrement" json:"id_user"`
	Username     string    `gorm:"not null" json:"username"`
	Name         string    `gorm:"not null" json:"name"`
	Mobile       string    `gorm:"not null" json:"mobile"`
	Email        string    `gorm:"not null" json:"email"`
	Password     string    `gorm:"not null" json:"password"`
	AccessLevel  string    `gorm:"null" json:"access_level"`
	Status       string    `gorm:"null" json:"status"`
	Img          string    `gorm:"default:-" json:"img"`
	DateCreated  time.Time `json:"date_created"`
	Timestamp    time.Time `json:"timestamp"`
}

type UserToken struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"not null" json:"email"`
	Token         string    `gorm:"not null" json:"token"`
	DateCreated  int64 `json:"date_created"`
}

type JSONData map[string]interface{}


func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.AccessLevel = "user"
	user.Status = "inactive"
	user.Img = "-"
	user.DateCreated = time.Now()
	user.Timestamp = time.Now()
	return nil
}

type UserRepository interface {
	RandomString(length int64) string
	CURLEmail(ctx context.Context, url string , user_token *UserToken, name string) (string, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Login(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, mobile string) (*User, error)
}

type UserUsecase interface {
	Create(ctx context.Context, req *req.UserCreateReq) (*User, error)
	Update(ctx context.Context, email string, req *req.UserPasswordReq) (*User, error)
	Login(ctx context.Context, req *req.UserLoginReq) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, mobile string) (*User, error)
}

func (User) TableName() string {
	return "user"
}

func (UserToken) TableName() string {
	return "user_token"
}
