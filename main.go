package main

import (
	routes "bookManagement/routes"
	"log"
	"net/http"
)

func main() {
	// api/v1/books endpoint
	injectedDB := "database.db"
	http.HandleFunc("/api/v1/books", func(w http.ResponseWriter, r *http.Request) {
		routes.AddBookHandler(w, r, injectedDB)
	})

	// api/v1/collection endpoints
	http.HandleFunc("/api/v1/collections", func(w http.ResponseWriter, r *http.Request) {
		routes.AddCollectionHandler(w, r, injectedDB)
	})

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
