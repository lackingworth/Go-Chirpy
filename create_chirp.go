package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"guthub.com/lackingworth/Go-Chirpy/internal/database"
)

var (
	BadWords = []string{"kerfuffle", "sharbert", "fornax"}
)

func badWordReplacer(txt string, badWords ...string) string {
	sliceOfStrings := strings.Fields(txt)
	for i, splitString := range sliceOfStrings {
		for _, badWord := range badWords {
			if strings.ToLower(splitString) == strings.ToLower(badWord) {
				sliceOfStrings[i] = "****"
			}
		}
	}
	return strings.Join(sliceOfStrings, " ")
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := badWordReplacer(params.Body, BadWords...)
	paramsCreateChirp := database.CreateChirpParams{Body: cleanedBody, UserID: params.UserID}

	var count int
	qrcErr := cfg.rawDbConn.QueryRowContext(r.Context(), "SELECT count(*) FROM users WHERE id=$1", paramsCreateChirp.UserID).Scan(&count)
	if qrcErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not check the associated user", qrcErr)
		return
	}
	if count == 0 {
		respondWithError(w, http.StatusBadRequest, "User does not exist", nil)
		return
	}

	chirp, dbErr := cfg.db.CreateChirp(r.Context(), paramsCreateChirp)
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", dbErr)
		return
	}

	returnedChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, returnedChirp)
}
