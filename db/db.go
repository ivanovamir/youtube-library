package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Book struct {
	ID   int64
	Name string
}

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	conn, err := sql.Open("postgres", "postgres://admin:admin@localhost:5432/library?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return &Database{db: conn}
}

func (p *Database) CreateBook(bookName string) error {
	_, err := p.db.Exec("INSERT INTO book (name) VALUES ($1)", bookName)
	return err
}

func (p *Database) GetAllBooks() ([]Book, error) {
	rows, err := p.db.Query("SELECT * FROM book")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var book Book
		if err = rows.Scan(&book.ID, &book.Name); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
