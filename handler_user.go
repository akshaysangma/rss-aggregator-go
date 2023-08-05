package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akshaysangma/rss-aggregator-go/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error Parsing Json : %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error Creating User : %v", err))
		return
	}

	respondWithJson(w, http.StatusCreated, user)
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, http.StatusOK, user)
}
