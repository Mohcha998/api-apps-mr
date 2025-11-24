package domain

type AdsPopup struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	ImageURL    string `json:"image_url" db:"image_url"`
	LinkURL     string `json:"link_url" db:"link_url"`
	FabImageURL string `json:"fab_image_url" db:"fab_image_url"`
	Status      int    `json:"status" db:"status"`
	DateCreated string `json:"date_created" db:"date_created"`
}

func (AdsPopup) TableName() string {
    return "ads_popup"
}