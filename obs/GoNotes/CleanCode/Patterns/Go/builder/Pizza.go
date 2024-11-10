package main

import (
	"fmt"
)

// Pizza представляет продукт, который мы будем строить.
type Pizza struct {
	size      string
	cheese    bool
	pepperoni bool
	veggies   bool
}

// PizzaBuilder интерфейс для строителей пиццы.
type PizzaBuilder interface {
	SetSize(size string) PizzaBuilder
	AddCheese() PizzaBuilder
	AddPepperoni() PizzaBuilder
	AddVeggies() PizzaBuilder
	Build() Pizza
}

// MargheritaPizzaBuilder конкретный строитель для пиццы Маргарита.
type MargheritaPizzaBuilder struct {
	pizza Pizza
}

func NewMargheritaPizzaBuilder() *MargheritaPizzaBuilder {
	return &MargheritaPizzaBuilder{}
}

func (b *MargheritaPizzaBuilder) SetSize(size string) PizzaBuilder {
	b.pizza.size = size
	return b
}

func (b *MargheritaPizzaBuilder) AddCheese() PizzaBuilder {
	b.pizza.cheese = true
	return b
}

func (b *MargheritaPizzaBuilder) AddPepperoni() PizzaBuilder {
	// Пицца Маргарита не содержит пепперони
	return b
}

func (b *MargheritaPizzaBuilder) AddVeggies() PizzaBuilder {
	b.pizza.veggies = true
	return b
}

func (b *MargheritaPizzaBuilder) Build() Pizza {
	return b.pizza
}

// PepperoniPizzaBuilder конкретный строитель для пиццы Пепперони.
type PepperoniPizzaBuilder struct {
	pizza Pizza
}

func NewPepperoniPizzaBuilder() *PepperoniPizzaBuilder {
	return &PepperoniPizzaBuilder{}
}

func (b *PepperoniPizzaBuilder) SetSize(size string) PizzaBuilder {
	b.pizza.size = size
	return b
}

func (b *PepperoniPizzaBuilder) AddCheese() PizzaBuilder {
	b.pizza.cheese = true
	return b
}

func (b *PepperoniPizzaBuilder) AddPepperoni() PizzaBuilder {
	b.pizza.pepperoni = true
	return b
}

func (b *PepperoniPizzaBuilder) AddVeggies() PizzaBuilder {
	b.pizza.veggies = true
	return b
}

func (b *PepperoniPizzaBuilder) Build() Pizza {
	return b.pizza
}

// Director управляет процессом строительства.
type Director struct {
	builder PizzaBuilder
}

// NewDirector создает нового директора с заданным строителем.
func NewDirector(builder PizzaBuilder) *Director {
	return &Director{builder: builder}
}

// ConstructPizza строит пиццу с заданными параметрами.
func (d *Director) ConstructPizza() {
	d.builder.SetSize("Large").AddCheese().AddVeggies()
	if _, ok := d.builder.(*PepperoniPizzaBuilder); ok {
		d.builder.AddPepperoni()
	}
}

func main() {
	// Строим пиццу Маргарита
	margheritaBuilder := NewMargheritaPizzaBuilder()
	director := NewDirector(margheritaBuilder)
	director.ConstructPizza()
	margheritaPizza := margheritaBuilder.Build()
	fmt.Println("Margherita Pizza:")
	fmt.Printf("Size: %s, Cheese: %t, Pepperoni: %t, Veggies: %t\n",
		margheritaPizza.size, margheritaPizza.cheese, margheritaPizza.pepperoni, margheritaPizza.veggies)

	// Строим пиццу Пепперони
	pepperoniBuilder := NewPepperoniPizzaBuilder()
	director = NewDirector(pepperoniBuilder)
	director.ConstructPizza()
	pepperoniPizza := pepperoniBuilder.Build()
	fmt.Println("Pepperoni Pizza:")
	fmt.Printf("Size: %s, Cheese: %t, Pepperoni: %t, Veggies: %t\n",
		pepperoniPizza.size, pepperoniPizza.cheese, pepperoniPizza.pepperoni, pepperoniPizza.veggies)
}
