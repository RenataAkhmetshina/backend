package handlers_test

import (
	"FlashcardLearningApp/user-service/db"
	"FlashcardLearningApp/user-service/handlers"
	"FlashcardLearningApp/user-service/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	err = gormDB.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return gormDB
}

func TestGetAllUsersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.GET("/users", handlers.GetAllUsers)

	testUsers := []models.User{
		{Username: "TestUserOne", Email: "testuser1@gmail.com", Bio: "TestUserBioOne", Password: "123456"},
		{Username: "TestUserTwo", Email: "testuser2@gmail.com", Bio: "TestUserBioTwo", Password: "123456"},
		{Username: "TestUserThree", Email: "testuser3@gmail.com", Bio: "TestUserBioThree", Password: "123456"},
	}

	if err := testDB.Create(&testUsers).Error; err != nil {
		t.Fatalf("failed to seed test data: %v", err)
	}

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseUsers []models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responseUsers); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(responseUsers) != len(testUsers) {
		t.Errorf("expected %d users, got %d", len(testUsers), len(responseUsers))
	}
}

func TestGetUserByIdHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.GET("/users/:id", handlers.GetUserById)

	existingUser := models.User{UserID: 5, Username: "TestUser", Email: "testuser@gmail.com", Bio: "TestUserBio"}
	testDB.Create(&existingUser)

	req, _ := http.NewRequest("GET", "/users/5", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var responseUser models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responseUser); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if responseUser.Username != "TestUser" {
		t.Errorf("expected username 'TestUser', got '%s'", responseUser.Username)
	}
}

func TestUpdateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.PUT("/users/:id", func(c *gin.Context) {
		c.Set("user_id", uint(10))
		c.Next()
	}, handlers.UpdateUser)

	originalUser := models.User{UserID: 10, Username: "TestUserOriginal", Email: "old@gmail.com", Bio: "OldBio"}
	testDB.Create(&originalUser)

	jsonBody := []byte(`{"email": "updated@gmail.com", "bio": "NewBio"}`)

	req, _ := http.NewRequest("PUT", "/users/10", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var dbUser models.User
	testDB.First(&dbUser, 10)
	if dbUser.Bio != "NewBio" || dbUser.Email != "updated@gmail.com" {
		t.Errorf("user was not updated correctly in DB: %+v", dbUser)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeFlashcardServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer fakeFlashcardServer.Close()

	t.Setenv("FLASHCARD_SERVICE_URL", fakeFlashcardServer.URL)
	t.Setenv("INTERNAL_SERVICE_KEY", "test-secret-key")

	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("user_id", uint(55))
		c.Next()
	}, handlers.DeleteUser)

	userToDelete := models.User{UserID: 55, Username: "TestUser", Email: "testuser@gmail.com"}
	testDB.Create(&userToDelete)

	req, _ := http.NewRequest("DELETE", "/users/55", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var dbUser models.User
	err := testDB.First(&dbUser, 55).Error
	if err == nil {
		t.Error("expected user to be deleted from database, but it still exists")
	}
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.POST("/auth/register", handlers.Register)

	jsonBody := []byte(`{"username": "NewUser", "email": "new@gmail.com", "bio": "Hello", "password": "password123"}`)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	var dbUser models.User
	if err := testDB.Where("username = ?", "NewUser").First(&dbUser).Error; err != nil {
		t.Fatalf("user was not saved in DB: %v", err)
	}
	if dbUser.Password == "password123" {
		t.Error("password was not hashed")
	}
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupTestDB(t)
	db.DB = testDB

	r := gin.New()
	r.POST("/auth/login", handlers.Login)

	r.POST("/auth/register", handlers.Register)
	regBody := []byte(`{"username": "TestUser", "email": "testuser@gmail.com", "password": "password123"}`)
	reqReg, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(regBody))
	reqReg.Header.Set("Content-Type", "application/json")
	rrReg := httptest.NewRecorder()
	r.ServeHTTP(rrReg, reqReg)

	loginBody := []byte(`{"username": "TestUser", "password": "password123"}`)
	reqLog, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
	reqLog.Header.Set("Content-Type", "application/json")

	rrLog := httptest.NewRecorder()
	r.ServeHTTP(rrLog, reqLog)

	if rrLog.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Response: %s", rrLog.Code, rrLog.Body.String())
	}

	var responseMap map[string]string
	json.Unmarshal(rrLog.Body.Bytes(), &responseMap)
	if responseMap["token"] == "" {
		t.Error("expected JWT token in response, got empty string")
	}
}
