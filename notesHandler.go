package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (connection *Connection) CreateNote(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var note Note
	note.Text = r.Form["text"][0]
	note.Username = r.Form["username"][0]
	note.When = time.Now()

	ret := connection.SaveNoteInDb(note)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if ret == true {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(note); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (connection *Connection) GetNotes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	jwtString := r.Header.Get("Authorization")
	log.Println("JWT: ", jwtString)
	username := GetUsernameFromToken(jwtString)
	if username != "" {
		// Token is OK. Fetch notes
		notes := connection.GetNotesByUser(username)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(notes); err != nil {
			panic(err)
		}
	} else {
		// Invalid token
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
	}
}
