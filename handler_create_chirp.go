package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"guthub.com/lackingworth/Go-Chirpy/internal/auth"
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

func validateChirp(body string) (string, error) {
	const maxChitpLength = 140
	if len(body) > maxChitpLength {
		return "", errors.New("Chirp is too long")
	}
	cleaned := badWordReplacer(body, BadWords...)
	return cleaned, nil
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	token, gbtErr := auth.GetBearerToken(r.Header)
	if gbtErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", gbtErr)
		return
	}

	userID, validationErr := auth.ValidateJWT(token, cfg.jwtSecret)
	if validationErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", validationErr)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	cleaned, validateErr := validateChirp(params.Body)
	if validateErr != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", validateErr)
		return
	}

	paramsCreateChirp := database.CreateChirpParams{Body: cleaned, UserID: userID}

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
