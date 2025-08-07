package req

import validation "github.com/go-ozzo/ozzo-validation"

type UserCreateReq struct {
	Name        string  `gorm:"not null" json:"name"`
	Mobile      string `gorm:"not null" json:"mobile"`
	Email       string  `gorm:"not null" json:"email"`
	Password    string  `gorm:"not null" json:"password"`
}

type UserLoginReq struct {
	Mobile      string `gorm:"not null" json:"mobile"`
	Email       string  `gorm:"not null" json:"email"`
	Password    string  `gorm:"not null" json:"password"`
}

type UserPasswordReq struct {
	Password string  `json:"password"`
}


func (r UserCreateReq) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Mobile, validation.Required),
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}


func (r UserPasswordReq) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Password, validation.Required),
	)
}

func (r UserLoginReq) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Password, validation.Required),
	)
}