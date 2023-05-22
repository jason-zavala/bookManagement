package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
		Status:  "error",
		Message: "Request to add book must include Author, and Title at a minimum.",
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

func TestGetBooksHandler(t *testing.T) {
	//setup
	databasePopulationHelper()
	defer cleanBooksTable()

	// Create a request
	req, err := http.NewRequest("GET", "/api/v1/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	r := httptest.NewRecorder()

	// Call the GetBooksHandler function with the request and response recorder
	GetBooksHandler(r, req, testDB)

	// Check the response status code
	if r.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.Code)
	}

	// Check the response body
	var books []Book
	err = json.Unmarshal(r.Body.Bytes(), &books)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the number of books retrieved, should be 25 based on the json file we read in the test helper method
	expectedCount := 25
	if len(books) != expectedCount {
		t.Errorf("Expected %d books, got %d", expectedCount, len(books))
	}
}

func TestBookFilterHandlerOneResult(t *testing.T) {
	databasePopulationHelper()
	defer cleanBooksTable()
	// Create a new request with the desired query parameters
	req, err := http.NewRequest("GET", "/api/v1/filter?title=1984", nil)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRecorder()
	FilterBooksHandler(r, req, testDB)

	if r.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, r.Code)
	}

	var book []Book
	err = json.Unmarshal(r.Body.Bytes(), &book)
	if err != nil {
		t.Fatal(err)
	}
	// Should only have returned 1 book
	if len(book) != 1 {
		t.Error("Expected one book got ", len(book))
	}
}

func TestBookFilterHandlerMultipleResults(t *testing.T) {
	databasePopulationHelper()
	defer cleanBooksTable()
	// Create a new request with the desired query parameters
	req, err := http.NewRequest("GET", "/api/v1/filter?genre=Fiction", nil)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRecorder()
	FilterBooksHandler(r, req, testDB)

	if r.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, r.Code)
	}

	var book []Book
	err = json.Unmarshal(r.Body.Bytes(), &book)
	if err != nil {
		t.Fatal(err)
	}
	// There should be 14 books here in the fiction genre
	if len(book) != 14 {
		t.Error("Expected one book got ", len(book))
	}
}

func TestBookFilterHandlerLongFilterURL(t *testing.T) {
	databasePopulationHelper()
	defer cleanBooksTable()
	// Create a new request with the desired query parameters
	req, err := http.NewRequest("GET", "/api/v1/filter?author=J.R.R.+Tolkien&genre=Fantasy", nil)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRecorder()
	FilterBooksHandler(r, req, testDB)

	if r.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, r.Code)
	}

	var books []Book
	err = json.Unmarshal(r.Body.Bytes(), &books)
	if err != nil {
		t.Fatal(err)
	}
	// There should be 2 books here in the fiction genre
	if len(books) != 2 {
		t.Error("Expected one book got ", len(books))
	}

	expectedTitles := []string{"The Lord of the Rings", "The Hobbit"}
	for i, book := range books {
		if book.Title != expectedTitles[i] {
			t.Errorf("Expected title '%s', got '%s'", expectedTitles[i], book.Title)
		}
	}
}

func databasePopulationHelper() {
	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Read books from JSON file
	books, err := readJSON("bookSetup.json")
	if err != nil {
		log.Fatal(err)
	}

	// Insert books into the database
	err = insertBookIntoDatabase(db, books)
	if err != nil {
		log.Println(err)
	}
}

func readJSON(filename string) ([]Book, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var books []Book

	err = json.Unmarshal(b, &books)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
		return nil, err
	}

	return books, nil
}

func insertBookIntoDatabase(db *sql.DB, books []Book) error {
	// Just making sure the db is empty so we can accurately get the count
	cleanBooksTable()

	query := "INSERT INTO Books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)"
	for _, book := range books {
		_, err := db.Exec(query, book.Title, book.Author, book.PublishedDate, book.Edition, book.Description, book.Genre)
		if err != nil {
			return fmt.Errorf("failed to insert book '%s': %v", book.Title, err)
		}
	}

	return nil
}

func cleanBooksTable() error {
	db, err := sql.Open("sqlite3", testDB)

	if err != nil {
		log.Fatal(err)
	}
	query := "DELETE FROM Books"

	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
