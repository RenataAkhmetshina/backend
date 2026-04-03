package models

type Category struct {
	CategoryID   uint   `gorm:"PrimaryKey" json:"category_id"`
	CategoryName string `gorm:"Unique" json:"category_name"`
}
