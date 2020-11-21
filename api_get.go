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
		data := getUserDisplayData(activeConnections.connections[sid.Value].User.Username)
		s, err = json.Marshal(data)
		break
	case "gain":
		data := getUserTransactionData(activeConnections.connections[sid.Value].User.Username, "gain")
		s, err = json.Marshal(data)
		break
	case "loss":
		data := getUserTransactionData(activeConnections.connections[sid.Value].User.Username, "loss")
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
	sid, ok := isLoggedIn(w, r)

	if !ok {
		fmt.Fprintf(w, "not logged")
		return
	}

	transactionType := mux.Vars(r)["type"]
	currencyType := mux.Vars(r)["currency"]

	if !isValidTransactionType(transactionType) || !isValidCurrecyType(currencyType) {
		fmt.Fprintf(w, "%s", string("[]"))
		return
	}

	var data []string

	if transactionType == "gain" {
		data = getAllCurrencies(currencyType)
	} else {
		data = getUserCurrencies(activeConnections.connections[sid.Value].User.Username, currencyType)
	}

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

	user := activeConnections.connections[sid.Value].User

	err = isValidTransaction(parameters, user.Username)
	if err != nil {
		replyReturnMessage(&w, "ERROR", "invalid transaction values")
		return
	}

	tData := createNewTransaction(parameters, user.Username)

	bData, err := getUserLatestPortfolio(user.Username)

	if err != nil {
		replyReturnMessage(&w, "ERROR", err.Error())
		return
	}

	err = updateBalanceWithTransaction(tData, bData)
	if err != nil {
		replyReturnMessage(&w, "ERROR", err.Error())
		return
	}

	refreshBalanceData(bData)

	err = insertUserPortfolio(user.Username, bData)
	if err != nil {
		replyReturnMessage(&w, "ERROR", "error inserting new portfolio balance")
		return
	}
	insertTransaction(tData)
	replyReturnMessage(&w, "OK", "transaction successful")
}

func apiV1Trade(w http.ResponseWriter, r *http.Request) {
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

	gain, loss := parseTradeParameters(parameters)

	user := activeConnections.connections[sid.Value].User

	err = isValidTransaction(gain, user.Username)
	if err != nil {
		replyReturnMessage(&w, "ERROR", "invalid trade values")
		return
	}
	err = isValidTransaction(loss, user.Username)
	if err != nil {
		replyReturnMessage(&w, "ERROR", "invalid trade values")
		return
	}

	tDataGain := createNewTransaction(gain, user.Username)
	tDataLoss := createNewTransaction(loss, user.Username)

	bData, err := getUserLatestPortfolio(user.Username)

	if err != nil {
		replyReturnMessage(&w, "ERROR", err.Error())
		return
	}

	err = updateBalanceWithTransaction(tDataGain, bData)
	if err != nil {
		replyReturnMessage(&w, "ERROR", err.Error())
		return
	}
	err = updateBalanceWithTransaction(tDataLoss, bData)
	if err != nil {
		replyReturnMessage(&w, "ERROR", err.Error())
		return
	}

	refreshBalanceData(bData)

	err = insertUserPortfolio(user.Username, bData)
	if err != nil {
		replyReturnMessage(&w, "ERROR", "error inserting new portfolio balance")
		return
	}
	insertTrade(tDataGain, tDataLoss)
	replyReturnMessage(&w, "OK", "trade successful")
}

func apiV1Username(w http.ResponseWriter, r *http.Request) {
	sid, ok := isLoggedIn(w, r)

	if !ok {
		fmt.Fprintf(w, "not logged")
		return
	}

	replyJSON, err := json.Marshal(activeConnections.connections[sid.Value].User.Username)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(w, string(replyJSON))
}

func apiV1Networth(w http.ResponseWriter, r *http.Request) {
	sid, ok := isLoggedIn(w, r)

	if !ok {
		fmt.Fprintf(w, "not logged")
		return
	}

	data := getUserNetworth(activeConnections.connections[sid.Value].User.Username)

	replyJSON, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(w, "%s", string(replyJSON))
}
