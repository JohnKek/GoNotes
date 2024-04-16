// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// !+broadcaster
type client struct {
	name string
	ch   chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			clients[cli] = true
			for c := range clients {
				cli.ch <- "User " + c.name + " is online"
			}

		case c := <-leaving:
			delete(clients, c)
			close(c.ch)
			for cli := range clients {
				cli.ch <- "User " + c.name + " has left"
			}
		}
	}
}

//!-broadcaster

// !+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	ch <- "Please enter your name: "
	input := bufio.NewScanner(conn)
	input.Scan()
	who := input.Text()
	ch <- "You are " + who
	messages <- who + " has arrived"
	cliet := client{name: who, ch: ch}
	entering <- cliet

	messageReceived := make(chan struct{}) // канал для сигнала о получении сообщения
	go func() {
		for input.Scan() {
			messages <- who + ": " + input.Text()
			messageReceived <- struct{}{} // отправляем сигнал о получении сообщения
		}
	}()
	idleTimer := time.NewTimer(1000 * time.Second)

	for {
		select {
		case <-idleTimer.C:
			leaving <- cliet
			messages <- who + " has left"
			err := conn.Close()
			if err != nil {
				return
			}
		case <-messageReceived:
			if !idleTimer.Stop() {
				<-idleTimer.C
			}
			idleTimer.Reset(1000 * time.Second)
		}
	}

}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

// !+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
