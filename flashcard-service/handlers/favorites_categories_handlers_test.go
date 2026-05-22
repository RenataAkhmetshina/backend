package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"FlashcardLearningApp/flashcard-service/db"
	"FlashcardLearningApp/flashcard-service/handlers"
	"FlashcardLearningApp/flashcard-service/models"

	"github.com/gin-gonic/gin"
)

func TestGetAllFavoriteCategoriesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.GET("/favorites", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	}, handlers.GetAllFavoriteCategories)

	testCategory := models.Category{CategoryID: 1, CategoryName: "TestCategory"}
	testDB.Create(&testCategory)

	testFavorites := []models.FavoriteCategories{
		{UserId: 1, CategoryId: 1},
	}
	testDB.Create(&testFavorites)

	req, _ := http.NewRequest("GET", "/favorites", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseFavorites []models.FavoriteCategories
	if err := json.Unmarshal(rr.Body.Bytes(), &responseFavorites); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(responseFavorites) != 1 {
		t.Errorf("expected 1 favorite category, got %d", len(responseFavorites))
	}
}

func TestAddToFavoriteCategoriesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.POST("/favorites/:id", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	}, handlers.AddToFavoriteCategories)

	testCategory := models.Category{CategoryID: 5, CategoryName: "TestCategory"}
	testDB.Create(&testCategory)

	req, _ := http.NewRequest("POST", "/favorites/5", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseMap map[string]string
	json.Unmarshal(rr.Body.Bytes(), &responseMap)
	if responseMap["message"] != "Added to favorites" {
		t.Errorf("expected success message, got: %s", responseMap["message"])
	}

	var dbFavorite models.FavoriteCategories
	if err := testDB.Where("user_id = ? AND category_id = ?", 1, 5).First(&dbFavorite).Error; err != nil {
		t.Errorf("Favorite record was not found in the test database: %v", err)
	}
}

func TestDeleteFromFavoriteCategoriesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.DELETE("/favorites/:id", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	}, handlers.DeleteFromFavoriteCategories)

	testCategory := models.Category{CategoryID: 10, CategoryName: "TestCategory"}
	testDB.Create(&testCategory)

	favoriteRecord := models.FavoriteCategories{UserId: 1, CategoryId: 10}
	testDB.Create(&favoriteRecord)

	req, _ := http.NewRequest("DELETE", "/favorites/10", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseMap map[string]string
	json.Unmarshal(rr.Body.Bytes(), &responseMap)
	if responseMap["message"] != "Successfully removed from favorites" {
		t.Errorf("expected success message, got: %s", responseMap["message"])
	}

	var dbFavorite models.FavoriteCategories
	err := testDB.Where("user_id = ? AND category_id = ?", 1, 10).First(&dbFavorite).Error
	if err == nil {
		t.Error("expected favorite record to be deleted, but it still exists in the database")
	}
}
