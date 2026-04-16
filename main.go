package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Set up routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/health", handleHealth)

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to debitask!\n")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK\n")
}
