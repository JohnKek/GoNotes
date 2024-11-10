package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {
	// Подключаемся к серверу
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error while connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	message := "Hello, world!"
	n := int32(10) // Количество пакетов для отправки
	buffer := make([]byte, 1024)
	// Отправляем количество пакетов в формате int32
	err = binary.Write(conn, binary.BigEndian, n)
	if err != nil {
		fmt.Println("Error while sending number:", err)
		return
	}
	/*	_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error while reading:", err)
			return
		}*/

	for i := int32(0); i < 15; i++ {
		// Отправляем сообщение серверу
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error while sending:", err)
			return
		}

		// Читаем ответ от сервера

		m, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error while reading:", err)
			return
		}

		fmt.Printf("Server response: %s\n", string(buffer[:m]))
	}
}
