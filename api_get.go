package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func apiV1Timeframe(w http.ResponseWriter, r *http.Request) {
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

func apiV1Table(w http.ResponseWriter, r *http.Request) {
	sid, ok := isLoggedIn(w, r)

	if !ok {
		fmt.Fprintf(w, "not logged")
		return
	}

	table := mux.Vars(r)["table"]

	if !isValidTable(table) {
		log.Println("invalid table")
		return
	}

	var s []byte
	var err error

	switch table {
	case "portfolio":
		data := getUserDisplayValues(activeConnections.connections[sid.Value].User.Username)
		s, err = json.Marshal(data)
		break
	}

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "%s", string(s))
}

func apiV1Currencies(w http.ResponseWriter, r *http.Request) {
	_, ok := isLoggedIn(w, r)

	if !ok {
		fmt.Fprintf(w, "not logged")
		return
	}

	currencyType := mux.Vars(r)["currency"]

	if !isValidCurrecyType(currencyType) {
		log.Println("invalid currency")
		return
	}

	data := getAllCurrencies(currencyType)
	s, err := json.Marshal(data)

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "%s", string(s))
}

func apiV1Transaction(w http.ResponseWriter, r *http.Request) {
	sid, ok := isLoggedIn(w, r)

	if !ok {
		fmt.Fprintf(w, "not logged")
		return
	}

	parameters := make(map[string]string)

	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type returnMessage struct {
		Status  string
		Message string
	}

	reply := &returnMessage{
		Status:  "OK",
		Message: "successfully inserted element",
	}

	user := activeConnections.connections[sid.Value].User
	tData, errMsg := isValidTransaction(parameters, user.Username)

	if len(errMsg) > 0 {
		reply.Status = "ERROR"
		reply.Message = errMsg
	} else {
		err := updateUserPortfolioTransaction(user.Username, tData)
		if err != nil {
			reply.Status = "ERROR"
			reply.Message = err.Error()
		}
	}

	replyJSON, err := json.Marshal(reply)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(w, string(replyJSON))
}
