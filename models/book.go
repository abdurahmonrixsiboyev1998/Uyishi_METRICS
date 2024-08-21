package models

import (
	"database/sql"
)

type Book struct {
	ID             int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"publishedDate"`
	ISBN          string `json:"isbn"`
}

func GetAllBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query("SELECT id, title, author, published_date, isbn FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.ISBN); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func CreateBook(db *sql.DB, book Book) (int, error) {
	stmt, err := db.Prepare("INSERT INTO books (title, author, published_date, isbn) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int
	if err := stmt.QueryRow(book.Title, book.Author, book.PublishedDate, book.ISBN).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func GetBook(db *sql.DB, id int) (*Book, error) {
	row := db.QueryRow("SELECT id, title, author, published_date, isbn FROM books WHERE id = $1", id)
	var book Book
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.ISBN); err != nil {
		return nil, err
	}
	return &book, nil
}

func UpdateBook(db *sql.DB, id int, book Book) error {
	_, err := db.Exec("UPDATE books SET title = $1, author = $2, published_date = $3, isbn = $4 WHERE id = $5", book.Title, book.Author, book.PublishedDate, book.ISBN, id)
	return err
}

func DeleteBook(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM books WHERE id = $1", id)
	return err
}