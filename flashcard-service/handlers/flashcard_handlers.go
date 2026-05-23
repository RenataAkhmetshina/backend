package handlers

import (
	"net/http"
	"strconv"

	"FlashcardLearningApp/flashcard-service/client"
	"FlashcardLearningApp/flashcard-service/db"
	"FlashcardLearningApp/flashcard-service/models"

	"github.com/gin-gonic/gin"
)

func GetAllFlashcards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	size := 5
	offset := (page - 1) * size

	query := db.DB.Model(&models.Flashcard{})

	title := c.Query("title")
	category_id, _ := strconv.Atoi(c.Query("category_id"))
	user_id, _ := strconv.Atoi(c.Query("user_id"))

	if title != "" {
		query = query.Where("title = ?", title)
	}

	if category_id != 0 {
		query = query.Where("category_id = ?", category_id)
	}

	if user_id != 0 {
		query = query.Where("user_id = ?", user_id)
	}

	var flashcards []models.Flashcard

	result := query.Limit(size).Offset(offset).Find(&flashcards)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching flashcards"})
		return
	}

	c.JSON(http.StatusOK, flashcards)
}

func GetFlashcardById(c *gin.Context) {
	var flashcard models.Flashcard
	flashcard_id := c.Param("id")

	if flashcard_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flashcard id is required"})
		return
	}

	result := db.DB.First(&flashcard, flashcard_id)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Flashcard not found"})
		return
	}

	c.JSON(http.StatusOK, flashcard)
}

func CreateFlashcard(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	token := c.GetHeader("Authorization")[7:]

	userClient := client.NewUserClient()
	exists, err := userClient.CheckUserExists(user_id.(uint), token)
	if err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}

	var newFlashcard models.Flashcard
	if err := c.ShouldBindJSON(&newFlashcard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newFlashcard.UserID = user_id.(uint)

	var category models.Category
	if err := db.DB.First(&category, newFlashcard.CategoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := db.DB.Create(&newFlashcard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Creation failed"})
		return
	}

	c.JSON(http.StatusCreated, newFlashcard)
}

func UpdateFlashcard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updatedData models.Flashcard
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var category models.Category
	if err := db.DB.First(&category, updatedData.CategoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	result := db.DB.Model(&models.Flashcard{}).
		Where("flashcard_id = ? AND user_id = ?", id, user_id.(uint)).
		Updates(updatedData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flashcard updated successfully"})
}

func DeleteFlashcard(c *gin.Context) {
	flashcardId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	userID, _ := c.Get("user_id")

	result := db.DB.Where("flashcard_id = ? AND user_id = ?", flashcardId, userID.(uint)).Delete(&models.Flashcard{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flashcard successfully deleted"})
}

func DeleteAllUserFlashcards(c *gin.Context) {
	userId := c.Param("userId")

	db.DB.Where("user_id = ?", userId).Delete(&models.Flashcard{})

	db.DB.Where("user_id = ?", userId).Delete(&models.FavoriteCategories{})

	c.JSON(http.StatusOK, gin.H{"message": "All user data wiped successfully"})
}
