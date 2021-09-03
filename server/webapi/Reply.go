package webapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type Reply struct {
	Status string `json:"status"`
	Data   string
}

func ReplyOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if data == nil {
		log.Println("data is nil")
		ReplyInternalError(w)
		return
	}

	reply := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: "OK",
		Data:   data,
	}

	b, err := json.Marshal(reply)
	if err != nil {
		log.Println(err)
		ReplyInternalError(w)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
		ReplyInternalError(w)
		return
	}
}

func ReplyInternalError(w http.ResponseWriter) {
	http.Error(w, "something went wrong :(", http.StatusInternalServerError)
}
