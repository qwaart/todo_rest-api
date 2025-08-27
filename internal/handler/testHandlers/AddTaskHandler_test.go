package test_Handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddTask(t *testing.T) {
	r := SetupRouter(t)
	//storage, _ := sqlite.New(":memory:")

	body := []byte(`{"title": "test task"}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}