package main

import (
	"net/http"

	"github.com/google/uuid"
	"guthub.com/lackingworth/Go-Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpsIDString := r.PathValue("chirpID")
	chirpID, parseErr := uuid.Parse(chirpsIDString)
	if parseErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", parseErr)
		return
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

	dbChirp, getChirpErr := cfg.db.GetChirp(r.Context(), chirpID)
	if getChirpErr != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", getChirpErr)
		return
	}
	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", getChirpErr)
		return
	}

	deleteErr := cfg.db.DeleteChirp(r.Context(), chirpID)
	if deleteErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
