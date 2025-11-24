package domain

type Version struct {
	ID          int64  `json:"id"`
	Img   		string `json:"img" db:"img"`
	Version 	string `json:"version" db:"version"`
	Feature     string `json:"feature" db:"feature"`
	CreatedDate string `json:"created_date" db:"created_date"`
}

func (Version) TableName() string {
    return "version"
}