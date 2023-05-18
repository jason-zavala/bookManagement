package main

import (
	books "bookManagement/routes"
	"log"
	"net/http"
)

func main() {
	// Register the handler for the /api/v1/books endpoint
	injectedDB := "database.db"
	http.HandleFunc("/api/v1/books", func(w http.ResponseWriter, r *http.Request) {
		books.AddBookHandler(w, r, injectedDB)
	})
	// Start the HTTP server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
