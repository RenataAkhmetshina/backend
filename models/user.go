package models

type User struct {
	UserID     uint        `gorm:"primaryKey" json:"user_id"`
	Username   string      `json:"username"`
	Email      string      `gorm:"Unique" json:"email"`
	Bio        string      `json:"bio"`
	Flashcards []Flashcard `gorm:"foreignKey:UserID" json:"flashcards"`
	Password   string      `json:"password"`
}
