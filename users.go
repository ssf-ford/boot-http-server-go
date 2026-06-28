package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"santnas/boot-http-server-course/internal/auth"
	"santnas/boot-http-server-course/internal/database"
)

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Email     string    `json:"email"`
}

func (c *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	user := userRequest{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	hashedPwd, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create hash", err)
		return
	}

	resData, err := c.db.CreateUser(r.Context(), database.CreateUserParams{Email: user.Email, HashedPassword: hashedPwd})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, userResponse{
		ID:        resData.ID,
		CreatedAt: resData.CreatedAt.Format(time.RFC3339),
		UpdatedAt: resData.UpdatedAt.Format(time.RFC3339),
		Email:     resData.Email,
	})
}

func (c *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	login := loginRequest{}
	err := decoder.Decode(&login)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	userDb, err := c.db.GetUserByEmail(r.Context(), login.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(login.Password, userDb.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, userResponse{
		ID:        userDb.ID,
		CreatedAt: userDb.CreatedAt.Format(time.RFC3339),
		UpdatedAt: userDb.UpdatedAt.Format(time.RFC3339),
		Email:     userDb.Email,
	})
}
