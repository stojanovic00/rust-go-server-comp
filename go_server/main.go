package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Define a handler function for handling HTTP requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Respond with a "200 OK" status
		w.WriteHeader(http.StatusOK)
		var mar = "MARKOOOOO"
		w.Write([]byte(mar))
		fmt.Fprint(w, "OK\n")
	})

	// Start the HTTP server on port 8002
	fmt.Println("Server listening on :8002")
	http.ListenAndServe(":8002", nil)
}
