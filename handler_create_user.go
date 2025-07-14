package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"guthub.com/lackingworth/Go-Chirpy/internal/auth"
	"guthub.com/lackingworth/Go-Chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	hashedPassword, hashErr := auth.HashPassword(params.Password)
	if hashErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", hashErr)
		return
	}

	dbCreateUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, dbErr := cfg.db.CreateUser(r.Context(), dbCreateUserParams)
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user", dbErr)
		return
	}
	returnedUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, returnedUser)
}
