package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, utcrErr := cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if utcrErr != nil {
		if errors.Is(utcrErr, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", utcrErr)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", utcrErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
