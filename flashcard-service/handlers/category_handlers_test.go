package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"FlashcardLearningApp/flashcard-service/db"
	"FlashcardLearningApp/flashcard-service/handlers"
	"FlashcardLearningApp/flashcard-service/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	err = gormDB.AutoMigrate(&models.Category{}, &models.Flashcard{}, &models.FavoriteCategories{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return gormDB
}

func TestGetAllCategoriesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.GET("/categories", handlers.GetAllCategories)

	testCategories := []models.Category{
		{CategoryName: "TestCategoryOne"},
		{CategoryName: "TestCategoryTwo"},
		{CategoryName: "TestCategoryThree"},
	}

	if err := testDB.Create(&testCategories).Error; err != nil {
		t.Fatalf("failed to seed test data: %v", err)
	}

	req, _ := http.NewRequest("GET", "/categories", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseCategories []models.Category
	if err := json.Unmarshal(rr.Body.Bytes(), &responseCategories); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedCount := len(testCategories)
	if len(responseCategories) != expectedCount {
		t.Errorf("expected %d categories, got %d", expectedCount, len(responseCategories))
	}

	if responseCategories[0].CategoryName != "TestCategoryOne" {
		t.Errorf("expected first category to be 'TestCategoryOne', got '%s'", responseCategories[0].CategoryName)
	}
}

func TestCreateCategoryHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testDB := setupTestDB(t)

	db.DB = testDB

	r := gin.New()
	r.POST("/categories", handlers.CreateCategory)

	jsonBody := []byte(`{"category_name": "TestCategory"}`)

	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	var responseMap map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &responseMap); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedMessage := "Category created successfully"
	if responseMap["message"] != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, responseMap["message"])
	}

	var dbCategory models.Category
	if err := testDB.Where("category_name = ?", "TestCategory").First(&dbCategory).Error; err != nil {
		t.Errorf("Category was not found in the test database: %v", err)
	}
}
