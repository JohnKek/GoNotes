package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// Создаем TCP-слушатель на порту 8080
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error while listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server started on :8080")

	for {
		// Принимаем входящее соединение
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	// Читаем количество пакетов из соединения
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error while reading:", err)
		return
	}

	var number int32
	err = binary.Read(bytes.NewReader(buffer[:n]), binary.BigEndian, &number)
	if err != nil {
		fmt.Println("Error while reading number:", err)
		return
	}
	_, err = conn.Write(buffer[:n])
	if err != nil {
		fmt.Println("Error while writing:", err)
		return
	}

	fmt.Printf("Number of packets to receive: %d\n", number)

	for i := 0; i < 10; i++ {
		// Читаем данные из соединения
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed by remote host")
				return
			}
			fmt.Println("Error while reading:", err)
			return
		}

		// Отправляем обратно то же сообщение
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error while writing:", err)
			return
		}
	}
}
