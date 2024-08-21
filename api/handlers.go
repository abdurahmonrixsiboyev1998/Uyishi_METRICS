package api

import (
	"database/sql"
	"encoding/json"
	"expvar"
	"metrics/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	requestCount = expvar.NewInt("requestCount")
	successCount = expvar.NewInt("successCount")
	errorCount   = expvar.NewInt("errorCount")
)

func GetBooksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		books, err := models.GetAllBooks(db)
		if err != nil {
			errorCount.Add(1)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
		successCount.Add(1)
	}
}

func CreateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			errorCount.Add(1)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := models.CreateBook(db, book)
		if err != nil {
			errorCount.Add(1)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": id})
		successCount.Add(1)
	}
}

func GetBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorCount.Add(1)
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		book, err := models.GetBook(db, id)
		if err != nil {
			errorCount.Add(1)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
		successCount.Add(1)
	}
}

func UpdateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorCount.Add(1)
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			errorCount.Add(1)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := models.UpdateBook(db, id, book); err != nil {
			errorCount.Add(1)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		successCount.Add(1)
	}
}

func DeleteBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorCount.Add(1)
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		if err := models.DeleteBook(db, id); err != nil {
			errorCount.Add(1)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		successCount.Add(1)
	}
}

func SetupRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/books", GetBooksHandler(db)).Methods("GET")
	r.HandleFunc("/books", CreateBookHandler(db)).Methods("POST")
	r.HandleFunc("/books/{id}", GetBookHandler(db)).Methods("GET")
	r.HandleFunc("/books/{id}", UpdateBookHandler(db)).Methods("PUT")
	r.HandleFunc("/books/{id}", DeleteBookHandler(db)).Methods("DELETE")
	r.Handle("/metrics", http.DefaultServeMux)
	return r
}
