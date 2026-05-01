package handlers

import (
	"net/http"
	"strconv"

	"FlashcardLearningApp/user-service/client"
	"FlashcardLearningApp/user-service/db"
	"FlashcardLearningApp/user-service/models"

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

	username := c.Query("username")
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

	authenticatedUserID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if id != strconv.FormatUint(uint64(authenticatedUserID.(uint)), 10) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own account"})
		return
	}

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
	userId := c.Param("id")

	authenticatedUserID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userId != strconv.FormatUint(uint64(authenticatedUserID.(uint)), 10) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account"})
		return
	}

	fcClient := client.NewFlashcardClient()

	if err := fcClient.WipeUserData(userId); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Could not sync with flashcard service"})
		return
	}

	result := db.DB.Delete(&models.User{}, userId)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User and data deleted successfully"})
}
