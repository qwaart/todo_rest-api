package handler_test_GetTaskByID

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"

	"rest_api/internal/db/sqlite"
	"rest_api/internal/handler"

	"github.com/gin-gonic/gin"
)

func setupRouter(t *testing.T) *gin.Engine {
	storage, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to init strage: %v", err)
	}
	
	r := gin.Default()
	taskHandler := handler.NewTaskHandler(storage)
	r.POST("/task", taskHandler.CreateTask)
	r.GET("/task/:id", taskHandler.GetTaskByID)

	return r
}

func TestGetTaskByID(t *testing.T) {
	r := setupRouter(t)

	// Add task
	taskJSON := `{"title": "Test task"}`
	req, _ := http.NewRequest("POST", "/task", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	//Decode response, for get id
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse resopnse: %v", err)
	}
	id := int(resp["id"].(float64)) // convert id from float64 to int

	//now i make a GET request for this id
	getReq, _ := http.NewRequest("GET", "/task/"+strconv.Itoa(id), nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, getReq)

	if w2.Code != http.StatusOK {
		t.Fatalf("expect status %d, got %d", http.StatusOK, w2.Code)
	}

	//Check that the data is correct
	var getResp map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if getResp["title"] != "Test task" {
		t.Errorf("expected task 'Test task', got %v", getResp["task"])
	}

	if getResp["completed"] != false {
		t.Errorf("expected completed=false, got %v", getResp["completed"])
	}
}