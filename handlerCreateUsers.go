package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	//Decoding the request body
	type parameter struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	errJSON := decoder.Decode(&params)
	if errJSON != nil {
		log.Println("could not access the request body")
		respondWithError(w, http.StatusBadRequest, "Wrong format for the request")
		return
	}
	// Generating a UUID
	userUUID := uuid.New()

	// Creating the User to the database

	timeNow := time.Now().UTC()
	var userDB database.CreateUserParams = database.CreateUserParams{
		ID:        userUUID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Name:      params.Name,
	}
	specificUser, err := cfg.DB.CreateUser(r.Context(), userDB)
	if err != nil {
		log.Println("Could not write the User to the database")
		respondWithError(w, http.StatusInternalServerError, "Error on the server-side, could not write to the database")
		return
	}
	respondUser := databaseUsertoUser(specificUser)
	respondWithJSON(w, http.StatusCreated, respondUser)

}
