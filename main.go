package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load("environment.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverPort := os.Getenv("PORT")
	connectionString := os.Getenv("CONNECTION_STRING")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Could not establish database connection")
	}
	dbQueries := database.New(db)

	config := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", readinessHandler)
	mux.HandleFunc("GET /v1/err", errHandler)
	mux.HandleFunc("POST /v1/users", config.CreateUserHandler)
	mux.HandleFunc("GET /v1/users", config.GetUserByApiKey)

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", serverPort)
	log.Fatal(server.ListenAndServe())

}
