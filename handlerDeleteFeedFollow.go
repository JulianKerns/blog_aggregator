package main

import (
	"log"
	"net/http"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) DeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDS := r.PathValue("feedFollowID")

	if feedFollowIDS == "" {
		log.Println("Bad Request")
		respondWithError(w, http.StatusBadRequest, "No FeedID provided in the Request")
		return
	}
	feedFollowID, err := uuid.Parse(feedFollowIDS)
	if err != nil {
		log.Println("Could not parse the FeedFollowID")
		respondWithError(w, http.StatusBadRequest, "Could not parse the FeedFollowID")
		return
	}

	errDelete := cfg.DB.DeleteFeedFollow(r.Context(), feedFollowID)
	if errDelete != nil {
		log.Println("Could not delete the record!")
		respondWithError(w, http.StatusInternalServerError, "Could not delete record from the Database")
		return
	}
	respondWithJSON(w, http.StatusNoContent, "")

}
