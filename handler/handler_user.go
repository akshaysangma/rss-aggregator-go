package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akshaysangma/rss-aggregator-go/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandleUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error Parsing Json : %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error Creating User : %v", err))
		return
	}

	RespondWithJson(w, http.StatusCreated, user)
}

func (apiCfg *ApiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	RespondWithJson(w, http.StatusOK, user)
}
