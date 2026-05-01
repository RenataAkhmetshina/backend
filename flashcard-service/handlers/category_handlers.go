package handlers

import (
	"FlashcardLearningApp/flashcard-service/db"
	"FlashcardLearningApp/flashcard-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllCategories(c *gin.Context) {
	var categories []models.Category
	result := db.DB.Find(&categories)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var newCategory models.Category

	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.DB.Create(&newCategory)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Error while creating a new category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}
