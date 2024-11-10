```go
package main

import "fmt"

// Database - интерфейс для баз данных
type Database interface {
	Connect() string
}

// MySQL - конкретная реализация базы данных MySQL
type MySQL struct{}

func (db *MySQL) Connect() string {
	return "Connected to MySQL database"
}

// PostgreSQL - конкретная реализация базы данных PostgreSQL
type PostgreSQL struct{}

func (db *PostgreSQL) Connect() string {
	return "Connected to PostgreSQL database"
}

// DatabaseFactory - интерфейс для фабрики баз данных
type DatabaseFactory interface {
	CreateDatabase() Database
}

// MySQLFactory - фабрика для создания MySQL
type MySQLFactory struct{}

func (f *MySQLFactory) CreateDatabase() Database {
	return &MySQL{}
}

// PostgreSQLFactory - фабрика для создания PostgreSQL
type PostgreSQLFactory struct{}

func (f *PostgreSQLFactory) CreateDatabase() Database {
	return &PostgreSQL{}
}

func main() {
	var factory DatabaseFactory

	// Создаем MySQL
	factory = &MySQLFactory{}
	mysql := factory.CreateDatabase()
	fmt.Println(mysql.Connect()) // Вывод: Connected to MySQL database

	// Создаем PostgreSQL
	factory = &PostgreSQLFactory{}
	postgresql := factory.CreateDatabase()
	fmt.Println(postgresql.Connect()) // Вывод: Connected to PostgreSQL database
}

```