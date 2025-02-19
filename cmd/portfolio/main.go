package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create a file server handler for the static directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Default handler for the homepage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body><h1>Welcome to my portfolio</h1><p>This is a simple portfolio site.</p></body></html>")
	})

	// Start server
	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
