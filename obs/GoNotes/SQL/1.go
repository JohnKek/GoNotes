package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // (1)
)

func main() {
	// подключиться к БД
	db, err := sql.Open("sqlite", "books.db") // (2)
	if err != nil {
		panic(err)
	}
	defer db.Close() // (3)

	err = db.Ping() // (4)
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ connected to books db")
	// ✓ connected to books db
}
