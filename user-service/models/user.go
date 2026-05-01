package models

type User struct {
	UserID   uint   `gorm:"primaryKey" json:"user_id"`
	Username string `json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Bio      string `json:"bio"`
	Password string `json:"-"`
}
