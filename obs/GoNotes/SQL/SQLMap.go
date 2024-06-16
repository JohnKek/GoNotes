package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // (1)
)

// начало решения
const (
	CREATE_QUERY = "create table if not exists map(key text primary key, val blob)"
	GET_QUERY    = "select val from map where key = ?"
	DELETE_QUERY = "delete from map where key = ?"
	SET_QUERY    = "insert into map(key, val) values (?, ?) on conflict (key) do update set val = excluded.val"
)

// SQLMap представляет карту, которая хранится в SQL-базе данных
type SQLMap struct {
	db *sql.DB
}

// NewSQLMap создает новую SQL-карту в указанной базе
func NewSQLMap(db *sql.DB) (*SQLMap, error) {
	_, err := db.Exec(CREATE_QUERY)
	return &SQLMap{db}, err
}

// Get возвращает значение для указанного ключа.
// Если такого ключа нет - возвращает ошибку sql.ErrNoRows.
func (m *SQLMap) Get(key string) (any, error) {
	var val any
	row := m.db.QueryRow(GET_QUERY, key)
	err := row.Scan(&val)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		panic(err)
	} else {
		return val, nil
	}
}

// Set устанавливает значение для указанного ключа.
// Если такой ключ уже есть - затирает старое значение (это не считается ошибкой).
func (m *SQLMap) Set(key string, val any) error {
	_, err := m.db.Exec(SET_QUERY, key, val)
	return err
}

// Delete удаляет запись карты с указанным ключом.
// Если такого ключа нет - ничего не делает (это не считается ошибкой).
func (m *SQLMap) Delete(key string) error {
	_, err := m.db.Exec(DELETE_QUERY, key)
	return err
}

// конец решения

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping() // (4)
	if err != nil {
		panic(err)
	}
	m, err := NewSQLMap(db)
	if err != nil {
		panic(err)
	}

	m.Set("name", "Alice")
	m.Set("age", 42)

	name, err := m.Get("name")
	fmt.Printf("name = %v, err = %v\n", name, err)
	// name = Alice, err = <nil>

	age, err := m.Get("age")
	fmt.Printf("age = %v, err = %v\n", age, err)
	// age = 42, err = <nil>

	m.Set("name", "Bob")
	name, err = m.Get("name")
	fmt.Printf("name = %v, err = %v\n", name, err)
	// name = Bob, err = <nil>

	m.Delete("name")
	name, err = m.Get("name")
	fmt.Printf("name = %v, err = %v\n", name, err)
	// name = <nil>, err = sql: no rows in result set
}
