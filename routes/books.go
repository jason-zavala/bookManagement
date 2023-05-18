package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	Edition       string `json:"edition"`
	Description   string `json:"description"`
	Genre         string `json:"genre"`
}

type Response struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	BookID string `json:"book_id,omitempty"`
}

func AddBookHandler(w http.ResponseWriter, r *http.Request, injectedDB string) {
	db, err := sql.Open("sqlite3", injectedDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Parse request body
	var book Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		// Return error response
		response := Response{
			Status: "error",
			Code:   "400",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Save the book to the database
	insertBookQuery := `
		INSERT INTO Books (title, author, published_date, edition, description, genre)
		VALUES (?, ?, ?, ?, ?, ?);`
	result, err := db.Exec(insertBookQuery, book.Title, book.Author, book.PublishedDate, book.Edition, book.Description, book.Genre)
	if err != nil {
		// Return error response
		response := Response{
			Status: "error",
			Code:   "500",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate the unique book ID
	bookID, _ := result.LastInsertId()

	// Return success response
	response := Response{
		Status: "success",
		Code:   "200",
		BookID: strconv.FormatInt(bookID, 10),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}