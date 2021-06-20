package webapi

import (
	"log"
	"net/http"

	"github.com/matejgrzinic/portfolio/appcontext"
)

func ApiBalance(appcontext *appcontext.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := appcontext.Portfolio.GetUserRefreshedBalance("Ace")
		if err != nil {
			log.Println(err)
		}

		ReplyOK(w, b)
	}
}
