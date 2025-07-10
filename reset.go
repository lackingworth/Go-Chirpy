package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "You do not have permission for this operation", nil)
		return
	}

	if err := cfg.db.Reset(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while deleting users", err)
	}

	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hits reset to 0 and database reset to initial state")); err != nil {
		log.Printf("Failed to write response from handlerReadiness: %v", err)
	}
}
