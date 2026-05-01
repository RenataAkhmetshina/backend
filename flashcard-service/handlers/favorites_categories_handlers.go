package handlers

import (
	"FlashcardLearningApp/flashcard-service/db"
	"FlashcardLearningApp/flashcard-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllFavoriteCategories(c *gin.Context) {
	var favorites []models.FavoriteCategories

	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	result := db.DB.Where("user_id = ?", user_id).Find(&favorites)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching favorites"})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func AddToFavoriteCategories(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	user_id, exists := c.Get("user_id")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Category ID"})
		return
	}

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	var category models.Category
	if err := db.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "The category you are trying to favorite does not exist"})
		return
	}

	favorite := models.FavoriteCategories{
		UserId:     user_id.(uint),
		CategoryId: uint(id),
	}

	if err := db.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This category is already in favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to favorites"})
}

func DeleteFromFavoriteCategories(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	user_id, _ := c.Get("user_id")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Category ID"})
		return
	}

	result := db.DB.Where("user_id = ? AND category_id = ?", user_id, id).Delete(&models.FavoriteCategories{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully removed from favorites"})
}
