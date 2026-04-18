package models

type FavoriteCategories struct {
	UserId     uint `gorm:"PrimaryKey" json:"user_id"`
	CategoryId uint `json:"category_id"`
}
