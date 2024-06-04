package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("ContentType", "application/json")

	dataJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error decoding the request parameters %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dataJSON)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type ErrorJSON struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, ErrorJSON{
		Error: msg,
	})

}
