package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"santnas/boot-http-server-course/internal/database"
)

type ChirpResult struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserID    string `json:"user_id"`
	Body      string `json:"body"`
}

func (c *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.db.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't list chirps", err)
		return
	}

	var chirpsResult []ChirpResult
	for _, chirp := range chirps {
		chirpsResult = append(chirpsResult, ChirpResult{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			UserID:    chirp.UserID.String(),
			Body:      sanitizeBody(chirp.Message),
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsResult)
}

func (c *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse chirp ID", err)
		return
	}

	chirp, err := c.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		//respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, ChirpResult{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		UserID:    chirp.UserID.String(),
		Body:      sanitizeBody(chirp.Message),
	})
}

func (c *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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

	respondWithJSON(w, http.StatusCreated, ChirpResult{
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
