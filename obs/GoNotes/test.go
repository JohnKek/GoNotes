package main

import (
	"fmt"
	"time"
)

func main() {
	// Получаем текущее время
	now := time.Now().UTC()
	fmt.Println(now)
	// Вычисляем завтрашнюю дату
	tomorrow := now.Add(24 * time.Hour)
	fmt.Println(tomorrow)
	// Устанавливаем время на 00:00
	tomorrowMidnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	fmt.Println(tomorrowMidnight)
}
