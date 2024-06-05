package main

import (
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) GetUserByApiKey(w http.ResponseWriter, r *http.Request) {

	apiKeySent := r.Header.Get("Authorization")

	apiKeyToCheck, ok := strings.CutPrefix(apiKeySent, "ApiKey ")
	if !ok {
		log.Println("Could not alter the string or the given string was empty")
		respondWithError(w, http.StatusUnauthorized, "wrong ApiKey, no authorization")
		return
	}
	specificUser, err := cfg.DB.GetUser(r.Context(), apiKeyToCheck)
	if err != nil {
		log.Println("Could not get the user from the database")
		respondWithError(w, http.StatusInternalServerError, "ApiKey not present or could not retrieve the data")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseUsertoUser(specificUser))

}
