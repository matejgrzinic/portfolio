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

	// fmt.Println(activeConnections.connections[sid.Value].User.Username, activeConnections.connections[sid.Value].User.Started)
	if !activeConnections.connections[sid.Value].User.Started {
		fmt.Fprintf(w, "first time")
		// err := templates.ExecuteTemplate(w, "index", data)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		return
	}

	data := getUserTimeframeData(activeConnections.connections[sid.Value].User.Username, "day")
	//var data graphData
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
