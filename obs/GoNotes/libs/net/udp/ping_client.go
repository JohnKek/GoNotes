package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Создаем UDP-адрес сервера
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// Создаем UDP-соединение
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	message := []byte("Hello, server!")

	// Отправляем сообщение серверу
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending:", err)
		return
	}

	// Читаем ответ от сервера
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	fmt.Printf("Server response: %s\n", string(buffer[:n]))
}
