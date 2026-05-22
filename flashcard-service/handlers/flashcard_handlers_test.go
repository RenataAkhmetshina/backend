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
)

func TestGetAllFlashcardsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.GET("/flashcards", handlers.GetAllFlashcards)

	testFlashcards := []models.Flashcard{
		{Title: "TestFlashcardOne", Image: "url", Text: "textOne", CategoryID: 1, UserID: 1},
		{Title: "TestFlashcardTwo", Image: "url", Text: "textTwo", CategoryID: 1, UserID: 1},
		{Title: "TestFlashcardThree", Image: "url", Text: "textThree", CategoryID: 1, UserID: 1},
	}

	/*
		type Flashcard struct {
			FlashcardID uint   `gorm:"PrimaryKey" json:"flashcard_id"`
			Title       string `json:"title"`
			Image       string `json:"image"`
			Text        string `json:"text"`
			CategoryID  uint   `gorm:"foreignKey:CategoryID" json:"category_id"`
			UserID      uint   `json:"user_id"`
		}
	*/

	if err := testDB.Create(&testFlashcards).Error; err != nil {
		t.Fatalf("failed to seed test data: %v", err)
	}

	req, _ := http.NewRequest("GET", "/flashcards", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseFlashcards []models.Flashcard
	if err := json.Unmarshal(rr.Body.Bytes(), &responseFlashcards); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedCount := len(testFlashcards)
	if len(responseFlashcards) != expectedCount {
		t.Errorf("expected %d flashcards, got %d", expectedCount, len(responseFlashcards))
	}

	if responseFlashcards[0].Title != "TestFlashcardOne" {
		t.Errorf("expected first flashcard to be 'TestFlashcardOne', got '%s'", responseFlashcards[0].Title)
	}
}

func TestGetFlashcardByIdHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.GET("/flashcards/:id", handlers.GetFlashcardById)

	testCategory := models.Category{CategoryID: 1, CategoryName: "TestCategory"}
	testDB.Create(&testCategory)

	testFlashcard := models.Flashcard{FlashcardID: 123, Title: "TestFlashcard", CategoryID: 1, UserID: 1}
	testDB.Create(&testFlashcard)

	req, _ := http.NewRequest("GET", "/flashcards/123", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseFlashcard models.Flashcard
	if err := json.Unmarshal(rr.Body.Bytes(), &responseFlashcard); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if responseFlashcard.Title != "TestFlashcard" {
		t.Errorf("expected title 'TestFlashcard', got '%s'", responseFlashcard.Title)
	}
}

func TestCreateFlashcardHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeUserServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"exists": true}`))
	}))
	defer fakeUserServer.Close()

	t.Setenv("USER_SERVICE_URL", fakeUserServer.URL)

	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.POST("/flashcards", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	}, handlers.CreateFlashcard)

	testCategory := models.Category{CategoryID: 1, CategoryName: "Languages"}
	testDB.Create(&testCategory)

	jsonBody := []byte(`{"title": "TestFlashcardOne", "image": "url", "text": "textOne", "category_id": 1}`)

	req, _ := http.NewRequest("POST", "/flashcards", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer fake-test-token-string-that-is-long-enough")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseFlashcard models.Flashcard
	if err := json.Unmarshal(rr.Body.Bytes(), &responseFlashcard); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if responseFlashcard.Title != "TestFlashcardOne" {
		t.Errorf("expected title 'TestFlashcardOne', got '%s'", responseFlashcard.Title)
	}

	var dbFlashcard models.Flashcard
	if err := testDB.Where("title = ?", "TestFlashcardOne").First(&dbFlashcard).Error; err != nil {
		t.Errorf("Flashcard was not found in the test database: %v", err)
	}
}

func TestUpdateFlashcardHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()

	r.PUT("/flashcards/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}, handlers.UpdateFlashcard)

	testCategory := models.Category{CategoryID: 1, CategoryName: "TestCategory"}
	testDB.Create(&testCategory)

	originalFlashcard := models.Flashcard{FlashcardID: 10, Title: "OldTitle", CategoryID: 1, UserID: 1}
	testDB.Create(&originalFlashcard)

	jsonBody := []byte(`{"title": "NewTitle", "category_id": 1}`)

	req, _ := http.NewRequest("PUT", "/flashcards/10", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseMap map[string]string
	json.Unmarshal(rr.Body.Bytes(), &responseMap)
	if responseMap["message"] != "Flashcard updated successfully" {
		t.Errorf("expected success message, got: %s", responseMap["message"])
	}

	var dbFlashcard models.Flashcard
	testDB.First(&dbFlashcard, 10)
	if dbFlashcard.Title != "NewTitle" {
		t.Errorf("expected title in DB to be updated to 'NewTitle', got '%s'", dbFlashcard.Title)
	}
}

func TestDeleteFlashcardHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()

	r.DELETE("/flashcards/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}, handlers.DeleteFlashcard)

	testCategory := models.Category{CategoryID: 1, CategoryName: "TestCategory"}
	testDB.Create(&testCategory)

	flashcardToDelete := models.Flashcard{FlashcardID: 55, Title: "To Be Deleted", CategoryID: 1, UserID: 1}
	testDB.Create(&flashcardToDelete)

	req, _ := http.NewRequest("DELETE", "/flashcards/55", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseMap map[string]string
	json.Unmarshal(rr.Body.Bytes(), &responseMap)
	if responseMap["message"] != "Flashcard successfully deleted" {
		t.Errorf("expected success message, got: %s", responseMap["message"])
	}

	var dbFlashcard models.Flashcard
	err := testDB.First(&dbFlashcard, 55).Error
	if err == nil {
		t.Error("expected flashcard to be deleted from database, but it still exists")
	}
}

func TestDeleteAllUserFlashcardsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.DELETE("/internal/flashcards/user/:userId", handlers.DeleteAllUserFlashcards)

	testCategory := models.Category{CategoryID: 1, CategoryName: "TestCategory"}
	testDB.Create(&testCategory)

	userFlashcards := []models.Flashcard{
		{FlashcardID: 1, Title: "FlashcardTtileOne", CategoryID: 1, UserID: 42},
		{FlashcardID: 2, Title: "FlashcardTtileTwo", CategoryID: 1, UserID: 42},
	}
	testDB.Create(&userFlashcards)

	otherFlashcard := models.Flashcard{FlashcardID: 3, Title: "OtherUsersFlashcard", CategoryID: 1, UserID: 99}
	testDB.Create(&otherFlashcard)

	req, _ := http.NewRequest("DELETE", "/internal/flashcards/user/42", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var count42 int64
	testDB.Model(&models.Flashcard{}).Where("user_id = ?", 42).Count(&count42)
	if count42 != 0 {
		t.Errorf("expected 0 flashcards left for user 42, got %d", count42)
	}

	var count99 int64
	testDB.Model(&models.Flashcard{}).Where("user_id = ?", 99).Count(&count99)
	if count99 != 1 {
		t.Errorf("expected flashcard of user 99 to be safe, but count is %d", count99)
	}
}
