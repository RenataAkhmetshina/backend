package main

import (
	"FlashcardLearningApp/flashcard-service/db"
	"FlashcardLearningApp/flashcard-service/handlers"
	"FlashcardLearningApp/flashcard-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	r := gin.Default()

	private := r.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		// Flashcards
		private.GET("/flashcards", handlers.GetAllFlashcards)
		private.GET("/flashcards/:id", handlers.GetFlashcardById)
		private.POST("/flashcards", handlers.CreateFlashcard)
		private.PUT("/flashcards/:id", handlers.UpdateFlashcard)
		private.DELETE("/flashcards/:id", handlers.DeleteFlashcard)

		// Categories
		private.GET("/categories", handlers.GetAllCategories)
		private.POST("/categories", handlers.CreateCategory)

		// FavoriteCategories
		private.GET("categories/favorites/", handlers.GetAllFavoriteCategories)
		private.PUT("/categories/:id/favorites", handlers.AddToFavoriteCategories)
		private.DELETE("/categories/:id/favorites", handlers.DeleteFromFavoriteCategories)
	}

	internal := r.Group("/internal")
	{
		internal.DELETE("/flashcards/user/:userId", handlers.DeleteAllUserFlashcards)
	}

	r.Run(":8082")
}
