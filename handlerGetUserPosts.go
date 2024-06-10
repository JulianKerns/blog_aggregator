package main

import (
	"log"
	"net/http"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
)

func (cfg *apiConfig) GetUserPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	const postLimit int = 20

	latestUserFeedPosts, err := cfg.DB.GetFeedPosts(r.Context(), database.GetFeedPostsParams{
		UserID: user.ID,
		Limit:  int32(postLimit),
	})
	if err != nil {
		log.Println("Could not retrieve the posts from the Database!")
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve the posts from the Database!")
		return
	}
	respondPost := make([]GetFeedPostsRow, len(latestUserFeedPosts))

	for i, latestPost := range latestUserFeedPosts {
		respondPost[i] = databaseFeedPosttoFeedPost(latestPost)

	}

	respondWithJSON(w, http.StatusOK, struct {
		Posts []GetFeedPostsRow `json:"posts"`
	}{
		Posts: respondPost,
	})

}
