package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
	mapMux := &sync.RWMutex{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn, repo, mapMux)
	}
}

func handleConnection(conn net.Conn, repo *Repo, mapMux *sync.RWMutex) {
	defer conn.Close()

	request, err := ParseRequest(conn)
	if err != nil {
		response := "HTTP/1.1 400 BadRequest\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
	}

	switch request.Method {
	case Get:
		handleGet(conn, repo, request, mapMux)
	case Put:
		handlePut(conn, repo, request, mapMux)
	default:
		response := "HTTP/1.1 501 NotImplemented\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
	}
}

func handleGet(conn net.Conn, repo *Repo, request *HttpRequest, mapMux *sync.RWMutex) {
	idStr := strings.TrimPrefix(request.Path, "/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response := "HTTP/1.1 400 BadRequest\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
		return
	}

	mapMux.RLock()
	entry, ok := repo.Entries[id]
	mapMux.RUnlock()
	if !ok {
		response := "HTTP/1.1 404 NotFound\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
		return
	}

	jsonBytes, err := json.Marshal(entry.Entity)
	if err != nil {
		response := "HTTP/1.1 500 NotImplemented\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
		return
	}

	jsonStr := string(jsonBytes)

	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", len(jsonStr), jsonStr)
	conn.Write([]byte(response))
}

func handlePut(conn net.Conn, repo *Repo, request *HttpRequest, mapMux *sync.RWMutex) {
	mapMux.RLock()
	_, exists := repo.Entries[request.Body.Id]
	mapMux.RUnlock()

	if exists {
		mapMux.RLock()
		entry, _ := repo.Entries[request.Body.Id]
		entry.Mux.Lock()
		repo.Entries[request.Body.Id] = *NewMapEntryWMux(*request.Body, entry.Mux)
		entry.Mux.Unlock()
		mapMux.RUnlock()
	} else {
		mapMux.Lock()
		repo.Entries[request.Body.Id] = *NewMapEntry(*request.Body)
		mapMux.Unlock()
	}

	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n")
	conn.Write([]byte(response))
}
