package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	sid, ok := isLoggedIn(w, r)

	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user := activeConnections.connections[sid.Value].User
	if !user.Started {
		fmt.Fprintf(w, "first time")
		return
	}

	type indexTemplate struct {
		Graph           graphData
		PortfolioValues []currencyData
	}

	data := &indexTemplate{
		Graph:           getUserTimeframeData(user.Username, "day"),
		PortfolioValues: getUserDisplayValues(user.Username),
	}

	err := templates.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if user, ok := authLogin(r.PostFormValue("uname"), r.PostFormValue("pwd")); ok {
			sid := generateSID()
			createCookie(w, sid)
			addConnection(sid, user)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	if _, ok := isLoggedIn(w, r); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "login", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	sid, ok := isLoggedIn(w, r)

	if ok {
		deleteCookie(w)
		deleteConnection(sid.Value)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if user, ok := createUser(r.PostFormValue("uname"), r.PostFormValue("pwd")); ok {
			insertNewUser(user)
			sid := generateSID()
			createCookie(w, sid)
			addConnection(sid, user)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	if _, ok := isLoggedIn(w, r); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "register", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
