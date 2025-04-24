package TestSingle

import (
	"fmt"
	"github.com/google/wire"
)

type Database struct {
	// поля вашей структуры
}

var ProviderSet = wire.NewSet(NewDatabase)

var instance *Database

func NewDatabase() *Database {
	if instance == nil {
		instance = &Database{}
	}
	return instance
}

func main() {
	db1 := NewDatabase()
	db2 := NewDatabase()

	fmt.Printf("db1: %p\n", db1)
	fmt.Printf("db2: %p\n", db2)

	if db1 == db2 {
		fmt.Println("db1 и db2 указывают на один и тот же экземпляр (синглтон)")
	} else {
		fmt.Println("db1 и db2 разные экземпляры")
	}
}
