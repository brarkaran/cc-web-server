package main

import (
	"bufio"
	"fmt"
	"os"
	"net"
	"path/filepath"
	"strings"
)

// handleConnection handles each client connection.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	requestLine, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	fields := strings.Fields(requestLine)
	if len(fields) < 3 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\nBad Request\r\n"))
		return
	}

	method, path := fields[0], fields[1]

	if method != "GET" {
		conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\nMethod Not Allowed\r\n"))
		return
	}

	// Serve the requested file
	serveFile(conn, path)
}

// serveFile serves a file or responds with a 404 error.
func serveFile(conn net.Conn, path string) {
	// Default to index.html for root
	if path == "/" {
		path = "/index.html"
	}

	filePath := filepath.Join("www", path)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\nNot Found\r\n"))
		return
	}

	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s", len(fileData), fileData)
	conn.Write([]byte(response))
}

// main sets up the server and listens for connections.
func main() {
	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()
	fmt.Println("Listening on port 80")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}

		// Handle each connection in a new goroutine for concurrency
		go handleConnection(conn)
	}
}
