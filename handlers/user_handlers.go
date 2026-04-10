package handlers

import (
	"net/http"
	"strconv"

	"FlashcardLearningApp/db"
	"FlashcardLearningApp/models"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	size := 5
	offset := (page - 1) * size

	query := db.DB.Model(&models.User{})

	username := c.Query("userName")
	email := c.Query("email")

	if username != "" {
		query = query.Where("username = ?", username)
	}

	if email != "" {
		query = query.Where("email = ?", email)
	}

	var users []models.User

	result := query.Limit(size).Offset(offset).Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserById(c *gin.Context) {
	var user models.User
	userId := c.Param("id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User id is required"})
		return
	}

	result := db.DB.First(&user, userId)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result := db.DB.Model(&models.User{}).Where("user_id = ?", id).Select("Bio", "Email").Updates(input)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	if err := db.DB.Where("user_id = ?", userId).Delete(&models.Flashcard{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete users's flashcards"})
	}

	result := db.DB.Delete(&models.User{}, userId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})

}
