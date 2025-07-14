package main

import (
	"net/http"
	"time"

	"guthub.com/lackingworth/Go-Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, gbtErr := auth.GetBearerToken(r.Header)
	if gbtErr != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", gbtErr)
		return
	}

	user, gufrtErr := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if gufrtErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", gufrtErr)
		return
	}

	accessToken, jwtErr := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if jwtErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", jwtErr)
		return
	}

	returnedResponse := response{
		Token: accessToken,
	}

	respondWithJSON(w, http.StatusOK, returnedResponse)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, gbtErr := auth.GetBearerToken(r.Header)
	if gbtErr != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", gbtErr)
		return
	}

	_, rrtErr := cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if rrtErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke session", rrtErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
