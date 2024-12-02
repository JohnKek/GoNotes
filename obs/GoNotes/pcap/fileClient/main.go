package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := "http://127.0.0.1:8367/dump-get"

	// Создаем тело запроса
	requestBody := map[string]string{
		"dumpName": "imsi_6565_6555_test_test_2024-11-21 15:06:10.605124389 +0300 MSK m=+10.695089228.pcap",
	}

	// Кодируем тело в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Ошибка при кодировании JSON: %v\n", err)
		return
	}

	// Отправляем POST-запрос
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Ошибка при отправке запроса: %v\n", err)
		return
	}
	defer resp.Body.Close() // Закрываем тело ответа после завершения работы

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка: сервер вернул статус %s\n", resp.Status)
		return
	}

	// Извлекаем имя файла из заголовка Content-Disposition
	contentDisposition := resp.Header.Get("Content-Disposition")
	var fileName string
	if contentDisposition != "" {
		// Извлекаем имя файла из заголовка
		parts := strings.Split(contentDisposition, ";")
		for _, part := range parts {
			if strings.HasPrefix(strings.TrimSpace(part), "filename=") {
				fileName = strings.Trim(strings.TrimPrefix(part, "filename="), "\"")
				break
			}
		}
	}

	// Если имя файла не найдено, используем имя по умолчанию
	if fileName == "" {
		fileName = "downloaded_file.pcap"
	}

	// Создаем файл для записи
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		return
	}
	defer outFile.Close() // Закрываем файл после завершения работы

	// Копируем содержимое ответа в файл
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}

	fmt.Printf("Файл успешно загружен и сохранен как '%s'\n", fileName)
}
