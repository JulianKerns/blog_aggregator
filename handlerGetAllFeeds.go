package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) GetAllFeedsHandler(w http.ResponseWriter, r *http.Request) {
	allFeedsDatabase, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		log.Println("Could not retrieve the feeds from the database")
		respondWithJSON(w, http.StatusInternalServerError, "Could not retrieve the feeds from the database")
		return
	}
	if len(allFeedsDatabase) == 0 {
		log.Println("No feeds present in the database")
		respondWithJSON(w, http.StatusInternalServerError, "No feeds present in the database")
		return
	}
	var allFeeds []Feed
	for _, feed := range allFeedsDatabase {
		changedFeed := databaseFeedtoFeed(feed)
		allFeeds = append(allFeeds, changedFeed)
	}

	respondWithJSON(w, http.StatusOK, allFeeds)

}
