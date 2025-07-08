package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hits reset to 0")); err != nil {
		log.Printf("Failed to write response from handlerReadiness: %v", err)
	}
}
