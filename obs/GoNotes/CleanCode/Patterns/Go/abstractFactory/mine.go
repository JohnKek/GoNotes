package main

import "fmt"

// Интерфейсы для продуктов
type Car interface {
	Drive() string
}

type Motorcycle interface {
	Ride() string
}

// Конкретные продукты
type Tesla struct {
	Model string
}

type BMW struct {
	Model string
}

type HarleyDavidson struct {
	Model string
}

type Yamaha struct {
	Model string
}

// Реализация методов для автомобилей
func (t *Tesla) Drive() string {
	return fmt.Sprintf("Driving a Tesla %s", t.Model)
}

func (b *BMW) Drive() string {
	return fmt.Sprintf("Driving a BMW %s", b.Model)
}

// Реализация методов для мотоциклов
func (h *HarleyDavidson) Ride() string {
	return fmt.Sprintf("Riding a Harley-Davidson %s", h.Model)
}

func (y *Yamaha) Ride() string {
	return fmt.Sprintf("Riding a Yamaha %s", y.Model)
}

// Интерфейс абстрактной фабрики
type VehicleFactory interface {
	CreateCar() Car
	CreateMotorcycle() Motorcycle
}

// Конкретные фабрики
type TeslaFactory struct{}

func (f *TeslaFactory) CreateCar() Car {
	return &Tesla{Model: "Model S"}
}

func (f *TeslaFactory) CreateMotorcycle() Motorcycle {
	return nil // Tesla не производит мотоциклы
}

type BMWFactory struct{}

func (f *BMWFactory) CreateCar() Car {
	return &BMW{Model: "X5"}
}

func (f *BMWFactory) CreateMotorcycle() Motorcycle {
	return nil // BMW не производит мотоциклы
}

type HarleyDavidsonFactory struct{}

func (f *HarleyDavidsonFactory) CreateCar() Car {
	return nil // Harley-Davidson не производит автомобили
}

func (f *HarleyDavidsonFactory) CreateMotorcycle() Motorcycle {
	return &HarleyDavidson{Model: "Street 750"}
}

type YamahaFactory struct{}

func (f *YamahaFactory) CreateCar() Car {
	return nil // Yamaha не производит автомобили
}

func (f *YamahaFactory) CreateMotorcycle() Motorcycle {
	return &Yamaha{Model: "YZF-R3"}
}

// Клиентский код
func clientCode(factory VehicleFactory) {
	if car := factory.CreateCar(); car != nil {
		fmt.Println(car.Drive())
	} else {
		fmt.Println("This factory does not produce cars.")
	}

	if motorcycle := factory.CreateMotorcycle(); motorcycle != nil {
		fmt.Println(motorcycle.Ride())
	} else {
		fmt.Println("This factory does not produce motorcycles.")
	}
}

func main() {
	var factory VehicleFactory

	// Используем фабрику Tesla
	factory = &TeslaFactory{}
	clientCode(factory)

	// Используем фабрику BMW
	factory = &BMWFactory{}
	clientCode(factory)

	// Используем фабрику Harley-Davidson
	factory = &HarleyDavidsonFactory{}
	clientCode(factory)

	// Используем фабрику Yamaha
	factory = &YamahaFactory{}
	clientCode(factory)
}
