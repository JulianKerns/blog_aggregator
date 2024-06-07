package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	//Decoding the request body
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
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
	feedFollowUUID := uuid.New()

	// Writing the FeedFollow to the database
	timeNow := time.Now().UTC()
	var feedFollowDB database.CreateFeedFollowParams = database.CreateFeedFollowParams{
		ID:        feedFollowUUID,
		UserID:    user.ID,
		FeedID:    params.FeedID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	specificFeedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollowDB)
	if err != nil {
		log.Println("Could not write the FeedFollow to the database")
		respondWithError(w, http.StatusInternalServerError, "Could not write to the database, FeedFollow may already exist")
		return
	}
	respondFeedFollow := databaseFeedFollowtoFeedFollow(specificFeedFollow)
	respondWithJSON(w, http.StatusCreated, respondFeedFollow)

}
