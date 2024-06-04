package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

//type apiConfig struct  {
//
//
//}

func main() {
	err := godotenv.Load("environment.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverPort := os.Getenv("PORT")

	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}

	mux.HandleFunc("GET /v1/healthz", readinessHandler)

	log.Printf("Serving on port: %s\n", serverPort)
	log.Fatal(server.ListenAndServe())

}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
