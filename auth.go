package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authData struct {
	Sid     string
	User    *userData
	Expires time.Time
}

type authDataMap struct {
	connections map[string]*authData
	mux         sync.Mutex
}

const cookieTimer time.Duration = 5

var activeConnections authDataMap

func authLogin(username string, password string) (*userData, bool) {
	user, err := getUserData(username)
	if err != nil {
		return user, false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Passwordhash), []byte(password)); err != nil {
		return user, false
	}

	return user, true
}

func createUser(username string, password string) (*userData, bool) { // change bool to error
	if len(username) == 0 || len(password) == 0 {
		return nil, false
	}

	users := getAllUsers()
	for _, u := range users {
		if username == u {
			return nil, false
		}
	}

	hashpw := encrypt(password)
	created := time.Now().Unix()

	return &userData{Username: username, Passwordhash: hashpw, Created: created, Active: false}, true
}

func createCookie(w http.ResponseWriter, sid string) {
	expiration := time.Now().Add(time.Minute * cookieTimer)
	cookie := http.Cookie{Name: "sid", Value: sid, Expires: expiration}
	http.SetCookie(w, &cookie)
}

func deleteCookie(w http.ResponseWriter) {
	expiration := time.Unix(0, 0)
	cookie := http.Cookie{Name: "sid", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func refreshCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Expires = time.Now().Add(time.Minute * cookieTimer)
	http.SetCookie(w, cookie)
	refreshConnectionExpiry(cookie.Value)
}

func isAuthValid(sid string) bool {
	activeConnections.mux.Lock()
	defer activeConnections.mux.Unlock()

	if v, ok := activeConnections.connections[sid]; ok {
		if v.Expires.Sub(time.Now()) > 0 {
			return true
		}
		delete(activeConnections.connections, sid)
	}
	return false
}

func refreshConnectionExpiry(sid string) {
	activeConnections.mux.Lock()
	defer activeConnections.mux.Unlock()

	activeConnections.connections[sid].Expires = time.Now().Add(time.Minute * cookieTimer)
}

func addConnection(sid string, user *userData) {
	activeConnections.mux.Lock()
	defer activeConnections.mux.Unlock()

	activeConnections.connections[sid] = &authData{Sid: sid, User: user, Expires: time.Now().Add(time.Minute * cookieTimer)}
}

func generateSID() string {
	for {
		randomString := strconv.FormatInt(rand.Int63(), 10)
		hash := encrypt(randomString)
		if _, ok := activeConnections.connections[hash]; !ok {
			return hash
		}
	}
}

func encrypt(s string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) (*http.Cookie, bool) {
	sid, err := r.Cookie("sid")
	if err != nil {
		if r.URL.String() != "/login" && r.URL.String() != "/register" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		// log.Println(err)
		return sid, false
	}

	if !isAuthValid(sid.Value) {
		deleteCookie(w)
		// log.Println("cookie is not valid")
		if r.URL.String() != "/login" && r.URL.String() != "/register" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		return sid, false
	}
	refreshCookie(w, sid)

	return sid, true
}
