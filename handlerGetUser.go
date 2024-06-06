package main

import (
	"net/http"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
)

func (cfg *apiConfig) GetUserByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJSON(w, http.StatusOK, databaseUsertoUser(user))

}
