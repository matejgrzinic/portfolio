package webapi

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matejgrzinic/portfolio/appcontext"
	"github.com/matejgrzinic/portfolio/portfolio"
)

func ApiTimeline(ctx appcontext.CTX) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("USER").(*portfolio.User)
		if !ok {
			log.Println("no user in context")
			ReplyInternalError(w)
			return
		}

		timeframe := mux.Vars(r)["timeframe"]
		b, err := ctx.Portfolio().UserTimeline(user, timeframe)
		if err != nil {
			log.Println(err)
			ReplyInternalError(w)
			return
		}

		ReplyOK(w, b)
	}
}
