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
	fmt.Println("✓ connected to db")
	// ✓ connected to books db
	query := `
    drop table if exists books;
    create table if not exists books(
        id integer primary key,
        title text,
        author text,
        num_pages integer,
        rating real
    );
`

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ created books table")
	// ✓ created books table
	query = `
    insert into books(title, author, num_pages, rating)
    values (?, ?, ?, ?)
`

	data := [][]any{
		{"The Catcher in the Rye", "J.D. Salinger", 277, 3.8},
		{"The Fellowship of the Ring", "J.R.R. Tolkin", 398, 4.36},
		{"The Giver", "Lois Lowry", 208, 4.13},
		{"The Da Vinci Code", "Dan Brown", 489, 3.84},
		{"The Alchemist", "Paulo Coelho", 197, 3.86},
	}

	for _, vals := range data {
		res, err := db.Exec(query, vals...) // (1)
		if err != nil {
			panic(err)
		}
		bookID, err := res.LastInsertId() // (2)
		fmt.Printf("added new book: id=%d, error=%v\n", bookID, err)
	}
	/*
	   added new book: id=1, error=<nil>
	   added new book: id=2, error=<nil>
	   added new book: id=3, error=<nil>
	   added new book: id=4, error=<nil>
	   added new book: id=5, error=<nil>
	*/

	query = "update books set author = ? where author = ?"
	res, err := db.Exec(query, "J.R.R. Tolkien", "J.R.R. Tolkin")
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	fmt.Printf("updated %d books, error=%v\n", count, err)
	// updated 1 books, error=<nil>

	query = "delete from books where rating < 4"
	res, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	count, err = res.RowsAffected()
	fmt.Printf("deleted %d books, error=%v\n", count, err)
	// deleted 3 books, error=<nil>

}
