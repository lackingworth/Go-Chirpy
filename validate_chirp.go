package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: badWordReplacer(params.Body, BadWords...),
	})
}
