package main

import (
	"encoding/json"
	"net/http"

	"guthub.com/lackingworth/Go-Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
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

	returnedResponse := response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	}
	respondWithJSON(w, http.StatusOK, returnedResponse)
}
