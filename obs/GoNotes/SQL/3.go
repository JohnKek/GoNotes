package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite" // (1)
)

func main() {
	db, err := sql.Open("sqlite", "books.db") // (2)
	if err != nil {
		panic(err)
	}
	defer db.Close() // (3)

	err = db.Ping() // (4)
	if err != nil {
		panic(err)
	}

	query := "select id, title, author, num_pages, rating from books"
	rows, err := db.Query(query) // (1)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // (2)

	for rows.Next() { // (3)
		var book Book
		err := rows.Scan(&book.id, &book.title, &book.author, &book.numPages, &book.rating) // (4)
		if err != nil {
			panic(err)
		}
		fmt.Println(book)
	}

	if err := rows.Err(); err != nil { // (5)
		panic(err)
	}
}

type Book struct {
	id       int
	title    string
	author   string
	numPages int
	rating   float64
}
