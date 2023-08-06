package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akshaysangma/rss-aggregator-go/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandleCreateFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		Feed_id uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON : %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.Feed_id,
		UserID:    u.ID,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed follow : %v", err))
		return
	}

	RespondWithJson(w, http.StatusCreated, feed)
}

func (apiCfg *ApiConfig) HandleGetFeedFollowsByUser(w http.ResponseWriter, r *http.Request, u database.User) {
	feeds, err := apiCfg.DB.GetFeedFollowsByUser(r.Context(), u.ID)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting feed follows for user : %v", err))
		return
	}

	if feeds == nil {
		feeds = []database.FeedFollow{}
	}

	RespondWithJson(w, http.StatusOK, feeds)
}

func (apiCfg *ApiConfig) HandleDeleteFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFollowIDVal := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDVal)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing feed follow ID : %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: u.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting feed follow : %v", err))
		return
	}

	RespondWithJson(w, http.StatusOK, struct{}{})
}
