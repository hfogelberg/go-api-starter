package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	mgo "gopkg.in/mgo.v2"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
)

type Note struct {
	Title string `json:"title"'`
	Text  string `json:"text"'`
}

type User struct {
	Email    string
	Username string
	Password []byte
}

type Notes []Note

func Signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := r.Form["email"][0]
	username := r.Form["username"][0]
	password := r.Form["password"][0]

	// Encrypt password
	pwd := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pwd, 10)
	log.Println("Hashed password", reflect.TypeOf(hashedPassword))
	if err != nil {
		panic(err)
	}

	fmt.Printf("The hashed password is : %s\n", string(hashedPassword))

	// Create user in Db
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()
	connection := session.DB("test").C("gousers")
	err2 := connection.Insert(&User{email, username, hashedPassword})
	if err2 != nil {
		panic(err2)
	}

	// Generat and return token
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login")
}

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
