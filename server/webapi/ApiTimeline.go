package webapi

import (
	"log"
	"net/http"

	"github.com/matejgrzinic/portfolio/appcontext"

	"github.com/gorilla/mux"
)

func ApiTimeline(appcontext *appcontext.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := "Ace"
		timeframe := mux.Vars(r)["timeframe"]

		data, err := appcontext.Portfolio.GetUserTimeline(user, timeframe) //appcontext.Db.Query.GetUserTimeline(user, timeframe)
		if err != nil {
			log.Println(err)
			ReplyInternalError(w)
			return
		}

		ReplyOK(w, data)
	}
}
