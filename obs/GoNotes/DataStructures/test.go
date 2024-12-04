package main

func main() {
	var ch chan int // Неинициализированный канал
	ch <- 1         // Попытка записи
}
