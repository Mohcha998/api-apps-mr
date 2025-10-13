package domain

import "time"

type QuotePool struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Text      string     `gorm:"column:text;type:text;not null" json:"text"`
	Author    *string    `gorm:"column:author;size:255" json:"author,omitempty"`
	CreatedBy *int64     `gorm:"column:created_by" json:"created_by,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
	IsActive  bool       `gorm:"column:is_active;default:1" json:"is_active"`
}

func (QuotePool) TableName() string {
	return "quotes_pool"
}

type QuoteAssigned struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"column:user_id;not null" json:"user_id"`
	QuotePoolID int64     `gorm:"column:quote_pool_id;not null" json:"quote_pool_id"`
	AssignDate  string    `gorm:"column:assign_date;type:date;not null" json:"assign_date"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (QuoteAssigned) TableName() string {
	return "quotes_assigned"
}

type QuotePoolHistory struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	QuotePoolID int64     `gorm:"column:quote_pool_id;not null" json:"quote_pool_id"`
	Text        string    `gorm:"column:text;type:text;not null" json:"text"`
	Author      *string   `gorm:"column:author;size:255" json:"author,omitempty"`
	ChangedBy   *int64    `gorm:"column:changed_by" json:"changed_by,omitempty"`
	ChangedAt   time.Time `gorm:"column:changed_at" json:"changed_at"`
}

func (QuotePoolHistory) TableName() string {
	return "quotes_pool_history"
}
