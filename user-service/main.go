package main

import (
	"FlashcardLearningApp/user-service/db"
	"FlashcardLearningApp/user-service/handlers"
	"FlashcardLearningApp/user-service/middleware"

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

	private := r.Group("/api/users")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/", handlers.GetAllUsers)
		private.GET("/:id", handlers.GetUserById)
		private.PUT("/:id", handlers.UpdateUser)
		private.DELETE("/:id", handlers.DeleteUser)
	}
	r.Run(":8081")
}
