package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"src/internal/handler"
	"src/internal/model"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testUserDB(t *testing.T) *gorm.DB {
	t.Helper()
	// Unique DSN so parallel tests / shared cache do not reuse the same DB.
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}))
	seed := []model.User{
		{ID: "1", Name: "John Doe", Email: "john.doe@example.com", Password: "secret"},
		{ID: "2", Name: "Jane Doe", Email: "jane.doe@example.com", Password: "secret"},
	}
	for i := range seed {
		require.NoError(t, db.Create(&seed[i]).Error)
	}
	return db
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := testUserDB(t)
	router := gin.New()
	userHandler := handler.NewUserHandler(db)
	router.POST("/users", userHandler.CreateUser)

	payload := map[string]string{
		"id":       "3",
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "secret",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code, "body: %s", rec.Body.String())

	var got map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

	assert.Equal(t, "3", got["id"])
	assert.Equal(t, "Alice", got["name"])
	assert.Equal(t, "alice@example.com", got["email"])
	assert.Equal(t, "secret", got["password"])
	assert.NotEmpty(t, got["created_at"])
	assert.NotEmpty(t, got["updated_at"])
}

func TestDeleteUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := testUserDB(t)
	router := gin.New()
	userHandler := handler.NewUserHandler(db)
	router.DELETE("/users/:id", userHandler.DeleteUser)
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "body: %s", rec.Body.String())

	var got map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "User deleted successfully", got["message"])
}

func TestUpdateUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := testUserDB(t)
	router := gin.New()
	userHandler := handler.NewUserHandler(db)
	router.PUT("/users/:id", userHandler.UpdateUser)
	payload := map[string]string{
		"id":       "1",
		"name":     "John Doe",
		"email":    "john.doe@example.com",
		"password": "secret",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "body: %s", rec.Body.String())

	var got map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "1", got["id"])
	assert.Equal(t, "John Doe", got["name"])
	assert.Equal(t, "john.doe@example.com", got["email"])
	assert.Equal(t, "secret", got["password"])
	assert.NotEmpty(t, got["created_at"])
	assert.NotEmpty(t, got["updated_at"])
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := testUserDB(t)
	router := gin.New()
	userHandler := handler.NewUserHandler(db)
	router.GET("/users", userHandler.GetUsers)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "body: %s", rec.Body.String())

	var got []map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

	assert.Len(t, got, 2)

	sort.Slice(got, func(i, j int) bool {
		idI, _ := got[i]["id"].(string)
		idJ, _ := got[j]["id"].(string)
		return strings.Compare(idI, idJ) < 0
	})

	assert.Equal(t, "1", got[0]["id"])
	assert.Equal(t, "John Doe", got[0]["name"])
	assert.Equal(t, "john.doe@example.com", got[0]["email"])

	assert.Equal(t, "2", got[1]["id"])
	assert.Equal(t, "Jane Doe", got[1]["name"])
	assert.Equal(t, "jane.doe@example.com", got[1]["email"])
}
