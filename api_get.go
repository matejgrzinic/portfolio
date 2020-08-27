package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func userTimeframeAPI(w http.ResponseWriter, r *http.Request) {
	for _, c := range r.Cookies() {
		if c.Name == "sid" && c.Path != "/" {
			c.Expires = time.Now().Add(-time.Minute)
		}
	}

	sid, ok := isLoggedIn(w, r)

	if !ok {
		fmt.Fprintf(w, "not logged")
		return
	}

	timeframe := mux.Vars(r)["timeframe"]

	if !isValidTimeframe(timeframe) {
		log.Println("invalid timeframe")
		return
	}

	data := getUserTimeframeData(activeConnections.connections[sid.Value].User.Username, timeframe)
	s, err := json.Marshal(data)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(w, "%s", string(s))

}
