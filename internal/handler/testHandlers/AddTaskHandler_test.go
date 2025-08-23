package handler_test_AddTask

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"rest_api/internal/db/sqlite"
	"rest_api/internal/handler"

	"github.com/gin-gonic/gin"
)

func setupRouter(storage *sqlite.Storage) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	taskHandler := handler.NewTaskHandler(storage)

	r.POST("/task", taskHandler.CreateTask)
	return r
}

func TestAddTask(t *testing.T) {
	storage, _ := sqlite.New(":memory:")
	router := setupRouter(storage)

	body := []byte(`{"title": "test task"}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}