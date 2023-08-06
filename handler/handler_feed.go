package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akshaysangma/rss-aggregator-go/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandleGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error could'nt get feeds : %v", err))
		return
	}

	if len(feeds) == 0 {
		feeds = []database.Feed{}
	}

	RespondWithJson(w, http.StatusOK, feeds)
}

func (apiCfg *ApiConfig) HandleCreateFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON : %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       params.URL,
		UserID:    u.ID,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed : %v", err))
		return
	}

	RespondWithJson(w, http.StatusCreated, feed)
}
