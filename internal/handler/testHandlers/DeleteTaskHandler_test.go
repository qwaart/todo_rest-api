package test_Handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"

	//"rest_api/internal/db/sqlite"
)

func TestDeleteTaskByID(t *testing.T) {
	r := SetupRouter(t)

	//Add task to be delete later
	taskJSON := `{"title":"Test task"}`
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
		t.Fatalf("failed to parse response: %v", err)
	}
	id := int(resp["id"].(float64))

	//DELETE task
	getReq, _ := http.NewRequest("DELETE", "/task/"+strconv.Itoa(id), nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, getReq)
	if w2.Code != http.StatusOK {
		t.Fatalf("expect status %d, got %d", http.StatusOK, w2.Code)
	}
}