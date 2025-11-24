package domain

import "time"

// -----------------------------
// Merchandise Tipe
// -----------------------------
type MerchandiseTipe struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"column:code;size:255" json:"code"`
	Name        string    `gorm:"column:name;size:255" json:"name"`
	Icon        string    `gorm:"column:icon;size:255" json:"icon"`
	Status      int       `gorm:"column:status" json:"status"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"`
}

func (MerchandiseTipe) TableName() string {
	return "merchandise_tipe"
}

// -----------------------------
// Merchandise Kategori
// -----------------------------
type MerchandiseKategori struct {
	ID                  int           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDMerchandiseTipe   int           `gorm:"column:id_merchandise_tipe;not null" json:"id_merchandise_tipe"`
	Code                string        `gorm:"column:code;size:255" json:"code"`
	Name                string        `gorm:"column:name;size:255" json:"name"`
	Status              int           `gorm:"column:status" json:"status"`
	CreatedDate         time.Time     `gorm:"column:created_date" json:"created_date"`
	NameMerchandiseTipe string        `gorm:"-" json:"name_merchandise_tipe"` // field tambahan untuk join
	Merchandise         []Merchandise `gorm:"-" json:"merchandise,omitempty"`
}

func (MerchandiseKategori) TableName() string {
	return "merchandise_kategori"
}

// -----------------------------
// Merchandise
// -----------------------------
type Merchandise struct {
	ID                    int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDMerchandiseTipe     int       `gorm:"column:id_merchandise_tipe;not null" json:"id_merchandise_tipe"`
	IDMerchandiseKategori int       `gorm:"column:id_merchandise_kategori;not null" json:"id_merchandise_kategori"`
	Code                  string    `gorm:"column:code;size:255" json:"code"`
	Judul                 string    `gorm:"column:judul;size:255" json:"judul"`
	Deskripsi             string    `gorm:"column:deskripsi;type:text" json:"deskripsi"`
	Gambar                string    `gorm:"column:gambar;size:255" json:"gambar"`
	RedirectLink          string    `gorm:"column:redirect_link;size:255" json:"redirect_link"`
	LinkNoPayment         string    `gorm:"column:link_no_payment;size:255" json:"link_no_payment"`
	CtaBtn                int       `gorm:"column:cta_btn" json:"cta_btn"`
	Status                int       `gorm:"column:status" json:"status"`
	CreatedDate           time.Time `gorm:"column:created_date" json:"created_date"`
}

func (Merchandise) TableName() string {
	return "merchandise"
}

// -----------------------------
// Aggregate Response
// -----------------------------
type MerchandiseAll struct {
	MerchandiseTipe []MerchandiseTipe        `json:"merchandise_tipe"`
	MerchandiseAll  []map[string]interface{} `json:"merchandise_all"`
	Merchandise     []Merchandise            `json:"merchandise"`
}
