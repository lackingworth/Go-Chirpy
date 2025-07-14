package main

import (
	"encoding/json"
	"net/http"

	"guthub.com/lackingworth/Go-Chirpy/internal/auth"
	"guthub.com/lackingworth/Go-Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, hashErr := auth.HashPassword(params.Password)
	if hashErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", hashErr)
		return
	}

	updateUserParams := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, updErr := cfg.db.UpdateUser(r.Context(), updateUserParams)
	if updErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", updErr)
		return
	}

	returnedResp := response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	}
	respondWithJSON(w, http.StatusOK, returnedResp)
}
