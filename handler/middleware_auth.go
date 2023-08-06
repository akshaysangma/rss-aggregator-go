package handler

import (
	"net/http"

	"github.com/akshaysangma/rss-aggregator-go/internal/auth"
	"github.com/akshaysangma/rss-aggregator-go/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig) MiddlewareAuth(Handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKeyFromHeader(r.Header)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid API Key")
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid API Key")
			return
		}
		Handler(w, r, user)
	}
}
