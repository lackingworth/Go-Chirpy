package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, dbErr := cfg.db.GetChirps(r.Context())
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't get chirps", dbErr)
		return
	}

	returnedSlice := make([]Chirp, 0, len(chirps))
	for _, chirp := range chirps {
		returnedInstance := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		returnedSlice = append(returnedSlice, returnedInstance)
	}
	respondWithJSON(w, http.StatusOK, returnedSlice)
}

func (cfg *apiConfig) handlerGetChirpsByID(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirpID", nil)
		return
	}

	chirp, dbErr := cfg.db.GetChirp(r.Context(), chirpID)
	if dbErr != nil {
		respondWithError(w, http.StatusNotFound, "Can't get chirp", dbErr)
		return
	}

	returnedChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, returnedChirp)
}
