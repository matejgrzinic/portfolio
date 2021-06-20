package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// appcontext := new(appcontext.AppContext)
	// appcontext.Ctx = context.TODO() // use it for queries??
	// appcontext.Db = db.NewDatabase("portfolio2")
	// appcontext.PriceData = external.NewPriceData(appcontext.Db.Db)
	// appcontext.Portfolio = portfolio.NewPortfolio(appcontext.Db, appcontext.PriceData)

	r := mux.NewRouter()
	//handler.SetupHandlers(r, appcontext)

	fmt.Printf("Starting server on: http://localhost:10000\n")

	s := &http.Server{
		Addr:         ":10000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
