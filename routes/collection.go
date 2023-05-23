package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Collection struct {
	CollectionID string `json:"collection_id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Books        []Book `json:"books"`
}

type CollectionResponse struct {
	CollectionID string `json:"collection_id,omitempty"`
	Message      string `json:"message,omitempty"`
	Status       string `json:"status"`
	Code         int    `json:"code"`
}

func AddCollectionHandler(w http.ResponseWriter, r *http.Request, injectedDB string) {
	db, err := sql.Open("sqlite3", injectedDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Parse request body
	var collection Collection
	err = json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		response := CollectionResponse{
			Status: "error",
			Code:   http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if collection.Description == "" || collection.Name == "" {
		response := CollectionResponse{
			Status:  "error",
			Message: "Collections must have at least a name and description.",
			Code:    http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	// Check if the collection already exists
	checkCollectionQuery := `SELECT collection_id FROM Collections WHERE name = ?;`
	var existingCollectionID int64
	err = db.QueryRow(checkCollectionQuery, collection.Name).Scan(&existingCollectionID)
	if err == nil {
		// Collection already exists, return existing collection ID
		response := CollectionResponse{
			Status:       "success",
			Code:         http.StatusOK,
			CollectionID: strconv.FormatInt(existingCollectionID, 10),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else if err != sql.ErrNoRows {
		// Error occurred during the database query
		response := CollectionResponse{
			Status: "error",
			Code:   http.StatusInternalServerError,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Save the collection to the database
	insertCollectionQuery := `INSERT INTO Collections (name, description) VALUES (?, ?);`

	result, err := db.Exec(insertCollectionQuery, collection.Name, collection.Description)
	if err != nil {
		response := CollectionResponse{
			Status: "error",
			Code:   http.StatusInternalServerError,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	collectionID, _ := result.LastInsertId()

	response := CollectionResponse{
		CollectionID: strconv.FormatInt(collectionID, 10),
		Status:       "success",
		Code:         http.StatusOK,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetCollectionsHandler(w http.ResponseWriter, r *http.Request, injectedDB string) {
	db, err := sql.Open("sqlite3", injectedDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the database to get all collections
	query := "SELECT collection_id, name, description FROM Collections"
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and create a list of collections
	var collections []Collection
	for rows.Next() {
		var collection Collection
		err := rows.Scan(&collection.CollectionID, &collection.Name, &collection.Description)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Query the database to get books associated with the collection
		bookQuery := "SELECT b.book_id, b.title, b.author FROM Books b INNER JOIN CollectionBooks cb ON b.book_id = cb.book_id WHERE cb.collection_id = ?"
		bookRows, err := db.Query(bookQuery, collection.CollectionID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer bookRows.Close()

		// Iterate over the book rows and create a list of books for the collection
		var books []Book
		for bookRows.Next() {
			var book Book
			err := bookRows.Scan(&book.BookID, &book.Title, &book.Author)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			books = append(books, book)
		}

		// Assign the books to the collection
		collection.Books = books

		collections = append(collections, collection)
	}

	// Check for any errors during row iteration
	err = rows.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode the collections list as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collections)
}

func AddBookToCollectionHandler(w http.ResponseWriter, r *http.Request, injectedDB string) {
	db, err := sql.Open("sqlite3", injectedDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var collectionToBookData struct {
		CollectionID string   `json:"collection_id"`
		BookIDs      []string `json:"book_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&collectionToBookData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the collection exists
	var existingCollectionID string
	err = db.QueryRow("SELECT collection_id FROM Collections WHERE collection_id = ?;", collectionToBookData.CollectionID).Scan(&existingCollectionID)
	if err == sql.ErrNoRows {
		response := Response{
			Status:  "error",
			Code:    http.StatusNotFound,
			Message: "Collection not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := Response{
			Status: "error",
			Code:   http.StatusInternalServerError,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if the books exist
	var existingBooks []string
	checkBooksQuery := fmt.Sprintf("SELECT book_id FROM Books WHERE book_id IN ('%s');", strings.Join(collectionToBookData.BookIDs, "','"))
	rows, err := db.Query(checkBooksQuery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var bookID string
		err := rows.Scan(&bookID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		existingBooks = append(existingBooks, bookID)
	}

	// Check if any books were not found
	missingBooks := make([]string, 0)
	for _, bookID := range collectionToBookData.BookIDs {
		found := false
		for _, existingBookID := range existingBooks {
			if bookID == existingBookID {
				found = true
				break
			}
		}
		if !found {
			missingBooks = append(missingBooks, bookID)
		}
	}

	if len(missingBooks) > 0 {
		response := Response{
			Status:  "error",
			Message: fmt.Sprintf("Books not found: %s", strings.Join(missingBooks, ", ")),
			Code:    http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Insert books into the collection
	insertQuery := `INSERT INTO CollectionBooks (collection_id, book_id) VALUES (?, ?);`
	for _, bookID := range collectionToBookData.BookIDs {
		_, err := db.Exec(insertQuery, collectionToBookData.CollectionID, bookID)
		if err != nil {
			response := Response{
				Status:  "error",
				Message: fmt.Sprintf("Failed to write %s into database", bookID),
				Code:    http.StatusInternalServerError,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response := Response{
		Status: "success",
		Code:   http.StatusOK,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
