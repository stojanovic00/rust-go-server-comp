package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type ShutdowServerFlag struct {
	sync.Mutex
	close bool
}

func NewShutdowServerFlag() *ShutdowServerFlag {
	return &ShutdowServerFlag{
		Mutex: sync.Mutex{},
		close: false,
	}
}

func main() {
	//Allocating resources
	repo := NewRepo()
	mapMux := &sync.RWMutex{}

	//Preparing thread pool
	poolSizeStr := os.Getenv("THREAD_POOL_SIZE")
	poolSize, err := strconv.Atoi(poolSizeStr)
	if err != nil {
		log.Fatal("Pool size not defined")
	}
	fmt.Printf("Thread pool size: %d\n", poolSize)

	connChan := make(chan net.Conn)
	wg := sync.WaitGroup{}

	//Starting thread pool
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func(threadNum int) {
			defer wg.Done()

			for {
				select {
				case conn, ok := <-connChan:
					{
						if !ok {
							return
						}
						fmt.Printf("Thread %d handles request\n", threadNum)
						handleConnection(conn, repo, mapMux)
					}
				}
			}
		}(i)
	}

	//Listening for requests
	tcp_address := os.Getenv("TCP_ADDRESS")
	listener, err := net.Listen("tcp", tcp_address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Listening for tcp requests on: %s\n", tcp_address)

	shutdownServerFlag := NewShutdowServerFlag()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err)
				continue
			}

			//Check for termination
			shutdownServerFlag.Lock()
			if shutdownServerFlag.close {
				close(connChan)
				shutdownServerFlag.Unlock()
				return
			} else {
				shutdownServerFlag.Unlock()
			}

			//Dispatch request to thread pool
			connChan <- conn
		}
	}()

	//Graceful shutdown
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	<-interruptCh
	shutdownServerFlag.Lock()
	shutdownServerFlag.close = true
	shutdownServerFlag.Unlock()

	sendShutdowRequest(tcp_address)

	fmt.Printf("Waiting for threads to close...\n")
	wg.Wait()
	fmt.Printf("Server shutdown\n")
}

// This is used to unstuck for loop  that waits for connection so it can terminate
func sendShutdowRequest(tcp_address string) {
	shutdownConn, err := net.Dial("tcp", tcp_address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		//TODO error handling
	}
	defer shutdownConn.Close()
	// Send an HTTP GET request
	request := "GET / HTTP/1.1\r\n\r\n"
	_, err = fmt.Fprintf(shutdownConn, request)
}

func handleConnection(conn net.Conn, repo *Repo, mapMux *sync.RWMutex) {
	defer conn.Close()

	request, err := ParseRequest(conn)
	if err != nil {
		response := "HTTP/1.1 400 BadRequest\r\nConnection: close\r\n\r\n" + err.Error()
		conn.Write([]byte(response))
		return
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

	//mapMux.RLock()
	entry, ok := repo.Entries.Load(id)
	//entry, ok := repo.Entries[id]
	//mapMux.RUnlock()
	if !ok {
		response := "HTTP/1.1 404 NotFound\r\nConnection: close\r\n\r\n"
		conn.Write([]byte(response))
		return
	}

	//jsonBytes, err := json.Marshal(entry.Entity)
	jsonBytes, err := json.Marshal(entry)
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
	//mapMux.RLock()
	//_, exists := repo.Entries[request.Body.Id]
	//mapMux.RUnlock()
	//
	//if exists {
	//	mapMux.RLock()
	//	entry, _ := repo.Entries[request.Body.Id]
	//	entry.Mux.Lock()
	//	repo.Entries[request.Body.Id] = *NewMapEntryWMux(*request.Body, entry.Mux)
	//	entry.Mux.Unlock()
	//	mapMux.RUnlock()
	//} else {
	//	mapMux.Lock()
	//	repo.Entries[request.Body.Id] = *NewMapEntry(*request.Body)
	//	mapMux.Unlock()
	//}

	repo.Entries.Swap(request.Body.Id, request.Body)

	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n")
	conn.Write([]byte(response))
}
