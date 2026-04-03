package main

import (
	"FlashcardLearningApp/db"
	"FlashcardLearningApp/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	r := gin.Default()

	// Flashcards
	r.GET("/flashcards", handlers.GetAllFlashcards)
	r.GET("/flashcards/:id", handlers.GetFlashcardById)
	r.POST("/flashcards", handlers.CreateFlashcard)
	r.PUT("/flashcards/:id", handlers.UpdateFlashcard)
	r.DELETE("/flashcards/:id", handlers.DeleteFlashcard)

	// Categpries
	r.GET("/categories", handlers.GetAllCategories)
	r.POST("/categories", handlers.CreateCategory)

	// Users
	r.GET("/users", handlers.GetAllUsers)
	r.GET("/users/:id", handlers.GetUserById)
	r.POST("/users", handlers.CreateUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	r.Run(":8080")
}
