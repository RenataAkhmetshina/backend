package main

import (
	"FlashcardLearningApp/db"
	"FlashcardLearningApp/handlers"
	"FlashcardLearningApp/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	r := gin.Default()

	public := r.Group("/auth")
	{
		public.POST("/register", handlers.Register)
		public.POST("/login", handlers.Login)
	}

	private := r.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		// Flashcards
		private.GET("/flashcards", handlers.GetAllFlashcards)
		private.GET("/flashcards/:id", handlers.GetFlashcardById)
		private.POST("/flashcards", handlers.CreateFlashcard)
		private.PUT("/flashcards/:id", handlers.UpdateFlashcard)
		private.DELETE("/flashcards/:id", handlers.DeleteFlashcard)

		// Categpries
		private.GET("/categories", handlers.GetAllCategories)
		private.POST("/categories", handlers.CreateCategory)

		// Users
		private.GET("/users", handlers.GetAllUsers)
		private.GET("/users/:id", handlers.GetUserById)
		private.PUT("/users/:id", handlers.UpdateUser)
		private.DELETE("/users/:id", handlers.DeleteUser)
	}

	r.Run(":8080")
}
