package main

import (
	"log"
	"net/http"

	auth "github.com/JulianKerns/blog_aggregator/internal/auth"
	database "github.com/JulianKerns/blog_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeyToCheck, err := auth.GetSentApiKey(r.Header)
		if err != nil {
			log.Printf("%s", err)
			respondWithError(w, http.StatusUnauthorized, "no authorization, false ApiKey")
			return
		}

		specificUser, errUser := cfg.DB.GetUser(r.Context(), apiKeyToCheck)

		if errUser != nil {
			log.Println("Could not get the user from the database")
			respondWithError(w, http.StatusInternalServerError, "ApiKey not present or could not retrieve the data")
			return
		}
		handler(w, r, specificUser)
	})
}
