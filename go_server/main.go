package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	tcp_address := os.Getenv("TCP_ADDRESS")
	listener, err := net.Listen("tcp", tcp_address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Listening for tcp requests on: %s\n", tcp_address)

	repo := NewRepo()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		handleConnection(conn, repo)
	}
}

func handleConnection(conn net.Conn, repo *Repo) {
	defer conn.Close()

	request, err := ParseRequest(conn)
	if err != nil {
		response := "HTTP/1.1 400 BadRequest\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
	}

	switch request.Method {
	case Get:
		handleGet(conn, repo, request)
	case Put:
		handlePut(conn, repo, request)
	default:
		response := "HTTP/1.1 501 NotImplemented\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
	}
}

func handleGet(conn net.Conn, repo *Repo, request *HttpRequest) {
	idStr := strings.TrimPrefix(request.Path, "/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response := "HTTP/1.1 400 BadRequest\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
	}

	entity, ok := repo.Entities[id]
	if !ok {
		response := "HTTP/1.1 404 NotFound\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
	}

	jsonBytes, err := json.Marshal(entity)
	if err != nil {
		response := "HTTP/1.1 500 NotImplemented\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
	}

	jsonStr := string(jsonBytes)

	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", len(jsonStr), jsonStr)
	conn.Write([]byte(response))
}

func handlePut(conn net.Conn, repo *Repo, request *HttpRequest) {
	repo.Entities[request.Body.Id] = *request.Body
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n")
	conn.Write([]byte(response))
}
