package handler

import (
	"context"
	"net/http"

	"github.com/matejgrzinic/portfolio/appcontext"
	"github.com/matejgrzinic/portfolio/portfolio"
)

type wrapper func(http.HandlerFunc, appcontext.CTX) http.HandlerFunc

func wrapperIsLoggedIn(f http.HandlerFunc, ctx appcontext.CTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.Clone(context.WithValue(r.Context(), "USER", &portfolio.User{Name: "Ace"}))

		//a, b := ctx.Currencies().GetCurrency("cryptocurrency", "BTC")
		//fmt.Println("test", a, b)
		// sid := ...
		// user = auth.SessionB
		f(w, r)
	}
}
