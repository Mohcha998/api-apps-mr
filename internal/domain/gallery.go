package domain

import "time"

type QuoteGallery struct {
	ID        int       `gorm:"column:id;primaryKey" json:"id"`
	Category  string    `gorm:"column:category" json:"category"`
	Images    string    `gorm:"column:images" json:"images"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (QuoteGallery) TableName() string {
	return "quote"
}
