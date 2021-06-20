package handler

import (
	"fmt"
	"net/http"

	"github.com/matejgrzinic/portfolio/appcontext"
)

func index(appcontext *appcontext.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "index")
	}
}
