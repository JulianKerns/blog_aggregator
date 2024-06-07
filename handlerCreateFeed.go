package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed_FeedFollow struct {
	Feed       Feed       `json:"feed"`
	FeedFollow FeedFollow `json:"feed_follow"`
}

func (cfg *apiConfig) CreateFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	//Decoding the request body
	type parameter struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	errJSON := decoder.Decode(&params)
	if errJSON != nil {
		log.Println("could not access the request body")
		respondWithError(w, http.StatusBadRequest, "Wrong format for the request")
		return
	}

	// Writing the Feed to the database
	feedUUID := uuid.New()
	timeNow := time.Now().UTC()
	var feedDB database.CreateFeedParams = database.CreateFeedParams{
		ID:        feedUUID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	}
	specificFeed, err := cfg.DB.CreateFeed(r.Context(), feedDB)
	if err != nil {
		log.Println("Could not write the Feed to the database")
		respondWithError(w, http.StatusInternalServerError, "Could not write to the database, feed may already exist")
		return
	}
	respondFeed := databaseFeedtoFeed(specificFeed)

	// Writing a FeedFollow to the created Feed by default
	feedFollowUUID := uuid.New()
	var feedFollowDB database.CreateFeedFollowParams = database.CreateFeedFollowParams{
		ID:        feedFollowUUID,
		UserID:    user.ID,
		FeedID:    feedUUID,
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

	respondWithJSON(w, http.StatusCreated, Feed_FeedFollow{
		Feed:       respondFeed,
		FeedFollow: respondFeedFollow,
	})

}
