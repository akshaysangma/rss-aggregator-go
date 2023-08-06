package handler

import "net/http"

func HandleReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJson(w, http.StatusOK, struct{}{})
}
