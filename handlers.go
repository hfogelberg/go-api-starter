package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetGreeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server says Hello!")
}
func GetNotes(w http.ResponseWriter, r *http.Request) {
	notes := Notes{
		Note{Title: "Hello 1", Text: "Lorem Ipsum 1"},
		Note{Title: "Hello 2", Text: "Lorem Ipsum 2"},
		Note{Title: "Hello 3", Text: "Lorem Ipsum 3"},
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(notes)
	if err != nil {
		panic(err)
	}
}

func GetNoteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteId := vars["id"]
	fmt.Fprintln(w, "Get note by id: ", noteId)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Note")
}
