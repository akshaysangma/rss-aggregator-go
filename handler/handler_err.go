package handler

import "net/http"

func HandleErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusNotFound, "Not found")
}
