package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// open a tcp listener on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	// close listener on exit
	defer ln.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		// accept connections
		conn, err := ln.Accept()

		if err != nil {
			// handle errors
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// handle connections in a separate routine
		go handleConnection(conn)
	}
}

// Connection Handler
func handleConnection(conn net.Conn) {

	// close connection on exit
	defer conn.Close()

	// notify console that a new connection has been made
	fmt.Printf("New connection from: %s\n", conn.RemoteAddr().String())

	// read data from the connection and print it to stdout
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("Received: %s\n", text)
	}

	// handle errors from reading data
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from connection:", err)
	}

}
