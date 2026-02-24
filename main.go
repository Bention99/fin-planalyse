package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve files under ./static at /
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	addr := ":8080"
	log.Printf("http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}