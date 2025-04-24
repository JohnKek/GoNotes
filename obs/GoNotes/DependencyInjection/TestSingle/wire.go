package TestSingle

import "github.com/google/wire"

func InitializeDatabase() *Database {
	wire.Build(ProviderSet)
	return &Database{} // Это не будет достигнуто, так как wire сгенерирует код для создания экземпляра
}
