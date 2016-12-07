package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
)

func handleNotes(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		createNote(w, r)
	case "GET":
		getNotes(w, r)
	default:
		http.Error(w, "Note supported", http.StatusMethodNotAllowed)
	}
}

func createNote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var note Note
	note.Text = r.Form["text"][0]
	note.Username = r.Form["username"][0]
	note.When = time.Now()

	log.Println("Note", note)

	// Hook up to Db
	db := context.Get(r, "database").(*mgo.Session)
	// insert it into the database
	if err := db.DB("test").C("gonotes").Insert(&note); err != nil {
		log.Println("Failed insert")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Insert OK")
	ret := Retval{
		Status:  100,
		Message: "OK",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func getNotes(w http.ResponseWriter, r *http.Request) {

}
