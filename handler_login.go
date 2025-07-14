package main

import (
	"encoding/json"
	"net/http"
	"time"

	"guthub.com/lackingworth/Go-Chirpy/internal/auth"
	"guthub.com/lackingworth/Go-Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	accessToken, jwtErr := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if jwtErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", jwtErr)
		return
	}

	refreshToken, mrtErr := auth.MakeRefreshToken()
	if mrtErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", mrtErr)
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	}

	_, crtErr := cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if crtErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", crtErr)
		return
	}

	returnedResponse := response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, returnedResponse)
}
