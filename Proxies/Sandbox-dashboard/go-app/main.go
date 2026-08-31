package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Use port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Register basic HTTP routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/hello", handleHello)

	fmt.Printf("Starting Go application on port %s...\n", port)
	
	// Listen on all network interfaces (0.0.0.0) so it's accessible from outside the container
	addr := ":" + port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

// Handler for the root URL "/"
func handleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the DevOps Learning Go App! The server is running successfully without Prometheus.\n"))
}

// Handler for the sample API URL "/api/hello"
func handleHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, DevOps Engineer! The application is responding correctly.\n"))
}
