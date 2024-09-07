package main

import (
	"encoding/json"
	"library/db"
	"net/http"
)

type CreateBook struct {
	Name string `json:"name"`
}

func main() {
	db := db.NewDatabase()
	r := http.NewServeMux()

	InitRoutes(r, db)
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func InitRoutes(r *http.ServeMux, pgsql *db.Database) {
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})

	r.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var newBook CreateBook
		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err = pgsql.CreateBook(newBook.Name); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		books, err := pgsql.GetAllBooks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		response, err := json.Marshal(books)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}
