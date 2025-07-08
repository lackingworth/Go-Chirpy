package main

import (
	"log"
	"net/http"
)

const (
	PORT = ":8080"
)

func main() {

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
