package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Collection struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CollectionResponse struct {
	CollectionID string `json:"collection_id,omitempty"`
	Message      string `json:"message,omitempty"`
	Status       string `json:"status"`
	Code         string `json:"code"`
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
			Code:   "400",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if collection.Description == "" || collection.Name == "" {
		response := CollectionResponse{
			Status:  "error",
			Message: "Collections must have at least a name and description.",
			Code:    "400",
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
			Code:         "200",
			CollectionID: strconv.FormatInt(existingCollectionID, 10),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else if err != sql.ErrNoRows {
		// Error occurred during the database query
		response := CollectionResponse{
			Status: "error",
			Code:   "500",
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
			Code:   "500",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	collectionID, _ := result.LastInsertId()

	response := CollectionResponse{
		CollectionID: strconv.FormatInt(collectionID, 10),
		Status:       "success",
		Code:         "200",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
