package models

type FavoriteCategories struct {
	UserId     uint `gorm:"primaryKey" json:"user_id"`
	CategoryId uint `gorm:"primaryKey" json:"category_id"`
}
