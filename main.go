package main

import (
	routes "bookManagement/routes"
	"log"
	"net/http"
)

func main() {
	// api/v1/books endpoint (this will handle both the get and the post methods)
	injectedDB := "routes/database.db"
	http.HandleFunc("/api/v1/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			routes.AddBookHandler(w, r, injectedDB)
		} else if r.Method == "GET" {
			routes.GetBooksHandler(w, r, injectedDB)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})

	// api/v1/collection endpoints
	http.HandleFunc("/api/v1/collections", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			routes.AddCollectionHandler(w, r, injectedDB)
		} else if r.Method == "GET" {
			routes.GetCollectionsHandler(w, r, injectedDB)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})

	//filter endpoint
	http.HandleFunc("/api/v1/filter", func(w http.ResponseWriter, r *http.Request) {
		routes.FilterBooksHandler(w, r, injectedDB)
	})

	//booksToCollection endpoint
	http.HandleFunc("/api/v1/booksToCollection", func(w http.ResponseWriter, r *http.Request) {
		routes.AddBookToCollectionHandler(w, r, injectedDB)
	})

	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
