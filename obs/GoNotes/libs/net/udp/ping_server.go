package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Создаем UDP-адрес
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// Создаем UDP-соединение
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP server started on :8080")

	buffer := make([]byte, 1024)

	for {
		// Читаем данные из соединения
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		fmt.Printf("Received message: %s from %s\n", string(buffer[:n]), remoteAddr)

		// Отправляем обратно то же сообщение
		_, err = conn.WriteToUDP(buffer[:n], remoteAddr)
		if err != nil {
			fmt.Println("Error writing:", err)
			continue
		}
	}
}
