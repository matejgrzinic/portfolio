package handler

import (
	"net/http"

	"github.com/matejgrzinic/portfolio/appcontext"
	"github.com/matejgrzinic/portfolio/webapi"

	"github.com/gorilla/mux"
)

func SetupHandlers(r *mux.Router, ctx appcontext.CTX) {
	r.Use(loggingMiddleware)

	r.HandleFunc("/api/balance", handleWithWrappers(webapi.ApiBalance(ctx), ctx,
		wrapperIsLoggedIn,
	))

	r.HandleFunc("/api/timeline/{timeframe}", handleWithWrappers(webapi.ApiTimeline(ctx), ctx,
		wrapperIsLoggedIn,
	))

	r.HandleFunc("/api/transaction/add", handleWithWrappers(webapi.ApiTransactionAdd(ctx), ctx,
		wrapperIsLoggedIn,
	))
}

func handleWithWrappers(f http.HandlerFunc, ctx appcontext.CTX, w ...wrapper) http.HandlerFunc {
	if len(w) == 0 {
		return f
	}
	return w[0](handleWithWrappers(f, ctx, w[1:cap(w)]...), ctx)
}
