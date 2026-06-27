package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"santnas/boot-http-server-course/internal/database"
)

func (c *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		UserID    string `json:"user_id"`
		Body      string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpReq := chirpRequest{}
	err := decoder.Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(chirpReq.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chripParams := database.CreateChirpParams{
		Message: chirpReq.Body,
		UserID:  chirpReq.UserID,
	}

	chirp, err := c.db.CreateChirp(r.Context(), chripParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		UserID:    chirp.UserID.String(),
		Body:      sanitizeBody(chirp.Message),
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
