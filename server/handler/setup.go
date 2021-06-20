package handler

import (
	"github.com/matejgrzinic/portfolio/appcontext"
	"github.com/matejgrzinic/portfolio/webapi"

	"github.com/gorilla/mux"
)

func SetupHandlers(r *mux.Router, appcontext *appcontext.AppContext) {
	r.Use(loggingMiddleware)
	r.HandleFunc("/", index(appcontext))
	r.HandleFunc("/api/balance", webapi.ApiTimeline(appcontext)).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/timeline/{timeframe}", webapi.ApiTimeline(appcontext)).Methods("GET", "OPTIONS")
}
