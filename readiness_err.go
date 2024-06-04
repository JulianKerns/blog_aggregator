package main

import (
	"net/http"
)

type ReadinessAndErr struct {
	Status string `json:"status"`
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, ReadinessAndErr{
		Status: "ok",
	})
}

func errHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
