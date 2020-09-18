package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var templates = template.Must(template.ParseGlob("templates/*"))
var db *mongo.Database

func main() {
	db = setup()

	go startUpdatePriceInterval()
	go startUpdateUserPortfolioInterval()

	activeConnections.connections = make(map[string]*authData)

	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/login", login)
	r.HandleFunc("/register", register)
	r.HandleFunc("/logout", logout)

	r.HandleFunc("/api/v1/timeframe/{timeframe}", apiV1Timeframe)
	r.HandleFunc("/api/v1/table/{table}", apiV1Table)
	r.HandleFunc("/api/v1/currencies/{currency}", apiV1Currencies)
	r.HandleFunc("/api/v1/transaction", apiV1Transaction)
	r.HandleFunc("/api/v1/username", apiV1Username)
	r.HandleFunc("/api/v1/networth", apiV1Networth)
	// change first two if you want to change how to access it over internet. last one is location on disk
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	fmt.Printf("Starting server on: http://localhost:10000\n")

	s := &http.Server{
		Addr:         ":10000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
