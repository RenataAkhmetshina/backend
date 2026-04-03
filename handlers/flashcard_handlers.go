package handlers

import (
	"net/http"
	"strconv"

	"FlashcardLearningApp/db"
	"FlashcardLearningApp/models"

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
	categoryId, _ := strconv.Atoi(c.Query("categoryId"))
	userId, _ := strconv.Atoi(c.Query("userId"))

	if title != "" {
		query = query.Where("title = ?", title)
	}

	if categoryId != 0 {
		query = query.Where("category_id = ?", categoryId)
	}

	if userId != 0 {
		query = query.Where("user_id = ?", userId)
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
	flashcardId := c.Param("id")

	if flashcardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flashcard id is required"})
		return
	}

	result := db.DB.First(&flashcard, flashcardId)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Flashcard not found"})
		return
	}

	c.JSON(http.StatusOK, flashcard)
}

func CreateFlashcard(c *gin.Context) {
	var newFlashcard models.Flashcard
	var category models.Category
	var user models.User

	if err := c.ShouldBindJSON(&newFlashcard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoryIdExists := db.DB.First(&category, newFlashcard.CategoryID)

	if categoryIdExists.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category with such id does not exist"})
		return
	}

	userIdExists := db.DB.First(&user, newFlashcard.UserID)

	if userIdExists.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User with such id does not exist"})
		return
	}

	result := db.DB.Create(&newFlashcard)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Error while creating a new flashcard"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Flashcard created successfully"})

}

func UpdateFlashcard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var updatedFlashcard models.Flashcard
	if err := c.ShouldBindJSON(&updatedFlashcard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var category models.Category
	var user models.User

	categoryIdExists := db.DB.First(&category, updatedFlashcard.CategoryID)

	if categoryIdExists.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category with such id does not exist"})
		return
	}

	userIdExists := db.DB.First(&user, updatedFlashcard.UserID)

	if userIdExists.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User with such id does not exist"})
		return
	}

	result := db.DB.Model(&models.Flashcard{}).Where("flashcard_id = ?", id).Updates(updatedFlashcard)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating the flashcard"})
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

	result := db.DB.Delete(&models.Flashcard{}, flashcardId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting the flashcard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flashcard successfully deleted"})

}
