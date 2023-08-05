package main

import "net/http"

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, struct{}{})
}
