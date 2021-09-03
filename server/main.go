package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/matejgrzinic/portfolio/appcontext"
	"github.com/matejgrzinic/portfolio/handler"
)

func SetupEnviormentVariables() {
	err := godotenv.Load("config.env")

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	SetupEnviormentVariables()

	ctx := appcontext.SetupAppContext()

	r := mux.NewRouter()
	handler.SetupHandlers(r, ctx)

	port := os.Getenv("PORT")
	fmt.Printf("Starting server on: http://localhost:%v\n", port)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
