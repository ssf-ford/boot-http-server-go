package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpRequest{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: sanitizeBody(chirp.Body),
	})
}

func sanitizeBody(body string) string {
	profanes := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	outputWords := []string{}

	for _, word := range words {
		if slices.Contains(profanes, strings.ToLower(word)) {
			outputWords = append(outputWords, "****")
		} else {
			outputWords = append(outputWords, word)
		}
	}

	return strings.Join(outputWords, " ")
}
