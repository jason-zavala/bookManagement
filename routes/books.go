package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	BookID        string `json:"book_id,omitempty"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	Edition       string `json:"edition"`
	Description   string `json:"description"`
	Genre         string `json:"genre"`
}

type Response struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	BookID  string `json:"book_id,omitempty"`
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
	//sanity checks here, things like making sure we at least have a title and author,

	if book.Author == "" || book.Title == "" {
		// Return error response
		response := Response{
			Status:  "error",
			Message: "Request to add book must include Author, and Title at a minimum.",
			Code:    "400",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	// Check if the book already exists
	checkBookQuery := `
		SELECT book_id FROM Books
		WHERE title = ? AND author = ? AND edition = ?;
	`
	var existingBookID int64
	err = db.QueryRow(checkBookQuery, book.Title, book.Author, book.Edition).Scan(&existingBookID)
	if err == nil {
		// Book already exists, return existing book ID
		response := Response{
			Status: "success",
			Code:   "200",
			BookID: strconv.FormatInt(existingBookID, 10),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else if err != sql.ErrNoRows {
		// Error occurred during the database query
		response := Response{
			Status: "error",
			Code:   "500",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Save the book to the database
	insertBookQuery := `
		INSERT INTO Books (title, author, published_date, edition, description, genre)
		VALUES (?, ?, ?, ?, ?, ?);
	`
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

func GetBooksHandler(w http.ResponseWriter, r *http.Request, injectedDB string) {
	db, err := sql.Open("sqlite3", injectedDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the database to get all books. Even tho we could use Select * notation here, we use the col names for clarity and readability
	query := "SELECT book_id, title, author, published_date, edition, description, genre FROM Books"
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and create a list of books
	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.PublishedDate, &book.Edition, &book.Description, &book.Genre)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// Check for any errors during row iteration
	err = rows.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode the books list as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}
