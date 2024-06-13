package main

import (
	"fmt"
)

// начало решения

// Map - карта "ключ-значение".
type Map[K comparable, V any] map[K]V

// Set устанавливает значение для ключа.
func (m Map[K, V]) Set(key K, val V) {
	m[key] = val
}

// Get возвращает значение по ключу.
func (m Map[K, V]) Get(key K) V {
	return m[key]
	// ...
}

// Keys возвращает срез ключей карты.
// Порядок ключей неважен, и не обязан совпадать
// с порядком значе ний из метода Values.
func (m Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
	// ...
}

// Values возвращает срез значений карты.
// Порядок значений неважен, и не обязан совпадать
// с порядком ключей из метода Keys.
func (m Map[K, V]) Values() []V {
	vals := make([]V, 0, len(m))
	for _, val := range m {
		vals = append(vals, val)
	}
	return vals
}

// конец решения

func main() {
	m := Map[string, int]{}
	m.Set("one", 1)
	m.Set("two", 2)

	fmt.Println(m.Get("one")) // 1
	fmt.Println(m.Get("two")) // 2

	fmt.Println(m.Keys())   // [one two]
	fmt.Println(m.Values()) // [1 2]
}
