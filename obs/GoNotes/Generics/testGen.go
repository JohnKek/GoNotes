package main

import "fmt"

// начало решения

// ZipMap возвращает карту, где ключи - элементы из keys, а значения - из vals.
func ZipMap[K comparable, V any](keys []K, vals []V) map[K]V {
	size := min(len(keys), len(vals))
	m := make(map[K]V, size)
	for i := range size {
		m[keys[i]] = vals[i]
	}
	return m
}

// конец решения

func main() {
	m1 := ZipMap([]string{"one", "two", "thr"}, []int{11, 22, 33})
	fmt.Println(m1)
	// map[one:11 two:22 thr:33]

	m2 := ZipMap([]string{"one"}, []int{11, 22, 33})
	fmt.Println(m2)
	// map[one:11]

	m3 := ZipMap([]string{"one", "two", "thr"}, []int{11})
	fmt.Println(m3)
	// map[one:11]

	m4 := ZipMap([]string{}, []int{11, 22, 33})
	fmt.Println(m4)
	// map[]
}
