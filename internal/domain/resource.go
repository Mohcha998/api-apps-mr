package domain

type Resource struct {
    ID            string `json:"id" db:"id"`
    Code          string `json:"code" db:"code"`
    Name          string `json:"name" db:"name"`
    Description   string `json:"description" db:"description"`
    Image         string `json:"image" db:"image"`
    Link          string `json:"link" db:"link"`
    TipeAsset     string `json:"tipe_asset" db:"tipe_asset"`
    JumlahHalaman string `json:"jumlah_halaman" db:"jumlah_halaman"`
    LinkHalaman   string `json:"link_halaman" db:"link_halaman"`
    Asset         string `json:"asset" db:"asset"`
    IsLockPremium string `json:"is_lock_premium" db:"is_lock_premium"`
    Status        string `json:"status" db:"status"`
    DateCreated   string `json:"date_created" db:"date_created"`
}


func (Resource) TableName() string {
    return "resource"
}
