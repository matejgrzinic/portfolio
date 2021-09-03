package webapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/matejgrzinic/portfolio/appcontext"
	"github.com/matejgrzinic/portfolio/portfolio"
)

func ApiTransactionAdd(ctx appcontext.CTX) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("USER").(*portfolio.User)
		if !ok {
			log.Println("no user in context")
			ReplyInternalError(w)
			return
		}

		var t portfolio.Transaction
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			log.Printf("decoding body into transaction: %v", err)
			ReplyInternalError(w)
			return
		}

		t.User = user.Name
		err = ctx.Portfolio().InsertTransaction(&t)
		if err != nil {
			log.Println(err)
			ReplyInternalError(w)
			return
		}

		ReplyOK(w, []interface{}{})
	}
}
