package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddCollectionHandlerSuccess(t *testing.T) {
	// Create a sample collection payload
	collection := Collection{
		Name:        "My Collection",
		Description: "A collection of my favorite books",
	}
	payload, _ := json.Marshal(collection)

	// Create a request with the sample payload
	req, err := http.NewRequest("POST", "/api/v1/collections", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	r := httptest.NewRecorder()

	// Call the AddCollectionHandler function with the request and response recorder
	AddCollectionHandler(r, req, testDB)

	// Check the response status code
	if r.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.Code)
	}

	// Check the response body
	var response CollectionResponse
	err = json.Unmarshal(r.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := CollectionResponse{
		CollectionID: "1", // Assuming the first collection has ID 1
		Status:       "success",
		Code:         "200",
	}

	if response.CollectionID != expectedResponse.CollectionID {
		t.Errorf("Expected collection ID %s, got %s", expectedResponse.CollectionID, response.CollectionID)
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %s, got %s", expectedResponse.Code, response.Code)
	}
}

func TestAddCollectionHandlerFail(t *testing.T) {
	// Create a sample collection payload with missing fields
	collection := Collection{}
	payload, _ := json.Marshal(collection)

	// Create a request with the sample payload
	req, err := http.NewRequest("POST", "/api/v1/collections", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRecorder()

	AddCollectionHandler(r, req, testDB)

	// Check the response status code
	if r.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, r.Code)
	}

	// Check the response body
	var response CollectionResponse
	err = json.Unmarshal(r.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := CollectionResponse{
		Status:  "error",
		Message: "Collections must have at least a name and description.",
		Code:    "400",
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %s, got %s", expectedResponse.Code, response.Code)
	}

	if response.Message != expectedResponse.Message {
		t.Errorf("Expected message %s, got %s", response.Message, expectedResponse.Message)
	}
}
