package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

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

	// Sanity checks, making sure we have at least a title and author
	if book.Author == "" || book.Title == "" {
		// Return error response
		response := Response{
			Status:  "error",
			Message: "Request to add book must include Author and Title at a minimum.",
			Code:    "400",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse the published date string
	publishedDate, err := parseDate(book.PublishedDate)
	if err != nil {
		// Return error response if the date parsing fails
		response := Response{
			Status:  "error",
			Message: "Failed to parse the published date. Valid formats for the date include YYYY, YYYY-MM, and YYYY-MM-DD",
			Code:    "400",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Save the book to the database
	insertBookQuery := "INSERT INTO Books (title, author, published_date, edition, description, genre)VALUES (?, ?, ?, ?, ?, ?);"
	result, err := db.Exec(insertBookQuery, book.Title, book.Author, publishedDate, book.Edition, book.Description, book.Genre)
	if err != nil {
		// Return error response
		response := Response{
			Status:  "error",
			Message: "Failed to save book to the database",
			Code:    "500",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// retrieve the unique book ID
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

func FilterBooksHandler(w http.ResponseWriter, r *http.Request, injectedDB string) {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Incorrectly formatted filter parameters",
			Code:    "400",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	//extract values from queryParams
	title := queryParams.Get("title")
	author := queryParams.Get("author")
	genre := queryParams.Get("genre")
	fromDate := queryParams.Get("fromData")
	toDate := queryParams.Get("toData")

	db, err := sql.Open("sqlite3", injectedDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "SELECT book_id, title, author, published_date, edition, description, genre FROM books WHERE 1=1"
	args := make([]interface{}, 0)

	if title != "" {
		query += " AND title = ?"
		args = append(args, title)
	}
	if author != "" {
		query += " AND author = ?"
		args = append(args, author)
	}
	if genre != "" {
		query += " AND genre = ?"
		args = append(args, genre)
	}
	if fromDate != "" {
		query += " AND published_date >= ?"
		args = append(args, fromDate)
	}
	if toDate != "" {
		query += " AND published_date <= ?"
		args = append(args, toDate)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	books := make([]Book, 0)

	for rows.Next() {
		var book Book
		err = rows.Scan(
			&book.BookID,
			&book.Title,
			&book.Author,
			&book.PublishedDate,
			&book.Edition,
			&book.Description,
			&book.Genre,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respJSON, err := json.Marshal(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

// This is a helper function to parse the dates because in testing the date wasnt being saved correctly in my DB
func parseDate(dateStr string) (time.Time, error) {
	if len(dateStr) == 4 {
		// The date string contains a year only
		return time.Parse("2006", dateStr)
	} else if len(dateStr) == 7 {
		// The date string contains a year and month
		return time.Parse("2006-01", dateStr)
	} else {
		// The date string contains a full date
		return time.Parse("2006-01-31", dateStr)
	}
}
