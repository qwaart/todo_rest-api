package test_Handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"
)

func TestUncompletedUpdateTaskByID(t *testing.T) {
	r := SetupRouter(t)

	//Add task
	taskJSON := `{"title": "Test task"}`
	req, _ := http.NewRequest("POST", "/task", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	//get id
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	id := int(resp["id"].(float64))

	getReq, _ := http.NewRequest("PATCH", "/task/"+strconv.Itoa(id)+"/uncompleted", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, getReq)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w2.Code)
	}

	var getResp map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to parse responce: %v", err)
	}

	if getResp["message"] != "task updated" {
		t.Errorf("expected message=task updated, got %v", getResp["message"])
	}
}