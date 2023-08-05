package main

import (
	"net/http"

	"github.com/akshaysangma/rss-aggregator-go/internal/auth"
	"github.com/akshaysangma/rss-aggregator-go/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKeyFromHeader(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid API Key")
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid API Key")
			return
		}
		handler(w, r, user)
	}
}
