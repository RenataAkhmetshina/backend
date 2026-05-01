package models

type Flashcard struct {
	FlashcardID uint   `gorm:"PrimaryKey" json:"flashcard_id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	Text        string `json:"text"`
	CategoryID  uint   `gorm:"foreignKey:CategoryID" json:"category_id"`
	UserID      uint   `json:"user_id"`
}
