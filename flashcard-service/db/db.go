package db

import (
	"log"

	"FlashcardLearningApp/flashcard-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=123 dbname=flashcard_service port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	log.Println("Connected successfully!")

	log.Println("Migration in progress...")
	err = DB.AutoMigrate(&models.Category{}, &models.Flashcard{}, &models.FavoriteCategories{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database migrated successfully!")
}
