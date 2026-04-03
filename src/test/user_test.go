package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"src/internal/handler"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	userHandler := handler.NewUserHandler()
	router.GET("/users", userHandler.GetUsers)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}

	if got[0]["id"] != "1" || got[0]["username"] != "John Doe" || got[0]["email"] != "john.doe@example.com" {
		t.Fatalf("unexpected first user payload: %+v", got[0])
	}

	if got[1]["id"] != "2" || got[1]["username"] != "Jane Doe" || got[1]["email"] != "jane.doe@example.com" {
		t.Fatalf("unexpected second user payload: %+v", got[1])
	}
}
