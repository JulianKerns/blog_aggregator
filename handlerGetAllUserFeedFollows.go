package main

import (
	"log"
	"net/http"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
)

func (cfg *apiConfig) GetAllUserFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	allUserFeedFollows, err := cfg.DB.GetAllUsersFeedFollows(r.Context(), user.ID)
	if err != nil {
		log.Println("Could not retrieve the feedfollows from the database")
		respondWithJSON(w, http.StatusInternalServerError, "Could not retrieve the feedfollows from the database")
		return
	}
	if len(allUserFeedFollows) == 0 {
		log.Println("No feedfollows present in the database")
		respondWithJSON(w, http.StatusInternalServerError, "No feedfollows present in the database")
		return
	}
	var allFeedFollows []FeedFollow
	for _, feedFollows := range allUserFeedFollows {
		changedFeedFollow := databaseFeedFollowtoFeedFollow(feedFollows)
		allFeedFollows = append(allFeedFollows, changedFeedFollow)
	}

	respondWithJSON(w, http.StatusOK, allFeedFollows)

}
