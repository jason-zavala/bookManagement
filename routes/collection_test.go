package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
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
		Status: "success",
		Code:   http.StatusOK,
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %d, got %d", expectedResponse.Code, response.Code)
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
		Code:    http.StatusBadRequest,
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %d, got %d", expectedResponse.Code, response.Code)
	}

	if response.Message != expectedResponse.Message {
		t.Errorf("Expected message %s, got %s", response.Message, expectedResponse.Message)
	}
}

func TestGetCollectionsHandler(t *testing.T) {
	//delete all collections so we dont affect our expected count for this test
	cleanCollectionsFromTestDatabase()
	collectionsToInsert := []Collection{
		{
			Name:        "Harry Potter",
			Description: "A Wizard saves the day!",
		},
		{
			Name:        "The Collected Sayings of MuadDib",
			Description: "Wisdom from the Lisan Al Gaib",
		},
	}

	for _, collect := range collectionsToInsert {
		err := insertCollection(collect)

		if err != nil {
			log.Fatal("Failed to insert into test database")
		}
	}

	// Create a new HTTP request to the collections endpoint
	req, err := http.NewRequest("GET", "/api/v1/collections", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	r := httptest.NewRecorder()

	// Call the GetCollectionsHandler function
	GetCollectionsHandler(r, req, testDB)

	// Check the response status code
	if r.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", r.Code, http.StatusOK)
	}

	// Parse the response body
	var collections []Collection
	err = json.Unmarshal(r.Body.Bytes(), &collections)
	if err != nil {
		t.Errorf("failed to parse response body: %v", err)
	}
	expectedLength := 2
	if len(collections) != expectedLength {
		t.Errorf("Expected length was %d, but expected %d", len(collections), expectedLength)
	}
}

func TestAddBookToCollectionHandlerSuccess(t *testing.T) {
	// Create a sample collection-to-book payload
	collectionToBookData := struct {
		CollectionID string        `json:"collection_id"`
		BookIDs      []interface{} `json:"book_ids"`
	}{
		CollectionID: "1",
		BookIDs:      []interface{}{"1"},
	}
	payload, _ := json.Marshal(collectionToBookData)

	cleanCollectionsFromTestDatabase()
	cleanBooksTable()

	c := Collection{
		Name:        "Harry Potter",
		Description: "A young boy engages in guerilla warfare in 90s London.",
	}

	insertCollection(c)

	b := Book{
		Title:  "Harry Potter and the Prisoner of Azkaban",
		Author: "J.K. Rowling",
	}

	insertBook(b)

	// Create a request with the sample payload
	req, err := http.NewRequest("POST", "/api/v1/bookToCollection", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	r := httptest.NewRecorder()

	// Call the AddBookToCollectionHandler function with the request and response recorder
	AddBookToCollectionHandler(r, req, testDB)

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
		Code:   http.StatusOK,
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %d, got %d", expectedResponse.Code, response.Code)
	}
}

func TestAddBookToCollectionHandlerCollectionNotFound(t *testing.T) {
	cleanBooksTable()

	b := Book{
		Title:  "Harry Potter and the Prisoner of Azkaban",
		Author: "J.K. Rowling",
	}

	insertBook(b)

	// payload with a non-existent collection
	collectionToBookData := struct {
		CollectionID string        `json:"collection_id"`
		BookIDs      []interface{} `json:"book_ids"`
	}{
		CollectionID: "999",
		BookIDs:      []interface{}{"1"},
	}
	payload, _ := json.Marshal(collectionToBookData)

	req, err := http.NewRequest("POST", "/api/v1/bookToCollection", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRecorder()

	AddBookToCollectionHandler(r, req, testDB)

	// Check the response status code
	if r.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, r.Code)
	}

	// Check the response body
	var response Response
	err = json.Unmarshal(r.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := Response{
		Status:  "error",
		Message: "Collection not found",
		Code:    http.StatusNotFound,
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %d, got %d", expectedResponse.Code, response.Code)
	}

	if response.Message != expectedResponse.Message {
		t.Errorf("Expected message %s, got %s", response.Message, expectedResponse.Message)
	}
}

func TestAddBookToCollectionHandlerBookNotFound(t *testing.T) {

	cleanCollectionsFromTestDatabase()

	c := Collection{
		Name:        "Harry Potter",
		Description: "A young boy engages in guerilla warfare in 90s London.",
	}

	insertCollection(c)

	// Create a payload with a non-existent book
	collectionToBookData := struct {
		CollectionID string        `json:"collection_id"`
		BookIDs      []interface{} `json:"book_ids"`
	}{
		CollectionID: "1",
		BookIDs:      []interface{}{"999"},
	}
	payload, _ := json.Marshal(collectionToBookData)

	req, err := http.NewRequest("POST", "/api/v1/bookToCollection", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRecorder()

	AddBookToCollectionHandler(r, req, testDB)

	if r.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, r.Code)
	}

	// Check the response body
	var response Response
	err = json.Unmarshal(r.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := Response{
		Status:  "error",
		Message: "Books not found: 999",
		Code:    http.StatusNotFound,
	}

	if response.Status != expectedResponse.Status {
		t.Errorf("Expected status %s, got %s", expectedResponse.Status, response.Status)
	}

	if response.Code != expectedResponse.Code {
		t.Errorf("Expected code %d, got %d", expectedResponse.Code, response.Code)
	}

	if response.Message != expectedResponse.Message {
		t.Errorf("Expected message %s, got %s", expectedResponse.Message, response.Message)
	}
}

func insertCollection(collection Collection) error {

	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "INSERT INTO Collections (name, description) VALUES (?, ?)"

	_, err = db.Exec(query, collection.Name, collection.Description)
	if err != nil {
		return err
	}

	return nil
}

func insertBook(book Book) error {
	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "INSERT INTO Books (title, author) VALUES (?, ?)"

	_, err = db.Exec(query, book.Title, book.Author)
	if err != nil {
		return err
	}

	return nil
}

func cleanCollectionsFromTestDatabase() error {
	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		log.Fatal(err)
	}

	query := "DELETE FROM Collections"

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
