package main

import (
	"encoding/json"
	"net/http"
	"time"

	"guthub.com/lackingworth/Go-Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, dbErr := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if dbErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	authErr := auth.VerifyPassword(params.Password, user.HashedPassword)
	if authErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", authErr)
		return
	}

	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	accessToken, jwtErr := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if jwtErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate access JWT", jwtErr)
		return
	}

	returnedResponse := response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: accessToken,
	}
	respondWithJSON(w, http.StatusOK, returnedResponse)
}
