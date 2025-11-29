package domain

type BirdtestUser struct {
	IDUser     int    `gorm:"column:id_user;primaryKey" json:"id_user"`
	Email      string `gorm:"column:email" json:"email"`
	IsBirdtest int    `gorm:"column:is_birdtest" json:"is_birdtest"`
}

func (BirdtestUser) TableName() string {
	return "user"
}

type BirdtestStatus struct {
	IsBirdtest int `json:"is_birdtest" gorm:"column:is_birdtest"`
}
