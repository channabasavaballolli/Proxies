package main

import (
	"net/http"          //used to construct HTTP request strcutures
	"net/http/httptest" //Go's mock library for web servers
	"testing"           //Go's standard package required for all unit tests
)

func TestGetTasksHandlerEmpty(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil) //creating a get request to /tasks
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder() //creates a response recorder to capture the output

	getTaskHandler(w, req) //call our handler directly

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	expectedBody := "[]\n"
	if w.Body.String() != expectedBody {
		t.Errorf("Eexpected body %q,got %q", expectedBody, w.Body.String())

	}

}
