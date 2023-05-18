package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testDB string = "testDB.db"

func TestAddBookHandlerSuccess(t *testing.T) {
	// Create a sample book payload
	book := Book{
		Title:         "Dune",
		Author:        "Frank Herbert",
		PublishedDate: "1965-08-01",
		Edition:       "1st Edition",
		Description:   "Paul Muad'Dib leads the Fremen on a conquest of revenge",
		Genre:         "Science Fiction",
	}
	payload, _ := json.Marshal(book)

	// Create a request with the sample payload
	req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	r := httptest.NewRecorder()

	// Call the AddBookHandler function with the request and response recorder
	AddBookHandler(r, req, testDB)

	// Check the response status code
	if r.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.Code)
	}

	// Check the response body
	var response Response
	err = json.Unmarshal(r.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := Response{
		Status: "success",
		Code:   "200",
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %s, got %s", expectedResponse.Code, response.Code)
	}
}

func TestAddBookHandlerFail(t *testing.T) {
	// Create a sample book payload
	book := Book{}
	payload, _ := json.Marshal(book)

	// Create a request with the sample payload
	req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRecorder()

	AddBookHandler(r, req, testDB)

	// Check the response status code
	if r.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, r.Code)
	}

	// Check the response body
	var response Response
	err = json.Unmarshal(r.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := Response{
		Status: "error",
		Code:   "400",
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %s, got %s", expectedResponse.Code, response.Code)
	}
}
