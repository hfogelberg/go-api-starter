package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func (connection *Connection) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Login")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]

	// Check if username is in Db
	user := User{}
	err = connection.Db.C("gousers").Find(bson.M{"username": username}).One(&user)

	log.Println("User is in Db", user.HashedPassword)
	if user.Username != "" {
		// Compare to password in Db
		pwd := []byte(password)
		err = bcrypt.CompareHashAndPassword(user.HashedPassword, pwd)
		fmt.Println(err) // nil means it is a match
		if err == nil {
			// Password OK. Generate token
			log.Println("Password is OK, Time to generate token")
			token := CreateToken(username)
			log.Println("We have a token!", token)

			ret := Retval{
				Status:  100,
				Token:   token,
				Message: "OK",
			}

			log.Println("Returning token")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ret); err != nil {
				panic(err)
			}

		} else {
			// Wrong password
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

	} else {
		// Wrong username
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
}

func (connection *Connection) Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var user User

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	password := r.Form["password"][0]
	pwd := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pwd, 10)
	user.Email = r.Form["email"][0]
	user.Username = r.Form["username"][0]
	user.HashedPassword = hashedPassword

	saveOk := connection.SaveUser(user)
	if saveOk == false {
		// todo return error
		log.Println("Error")
	}

	token := CreateToken(user.Username)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		panic(err)
	}
}

func (connection *Connection) CreateNote(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	err = connection.Db.C("gonotes").Insert(&note)
	if err != nil {
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

func (connection *Connection) GetNotes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	jwtString := r.Header.Get("Authorization")
	log.Println("JWT: ", jwtString)
	tokenIsValid := tokenIsValid(jwtString)
	log.Println("tokenIsValid: ", tokenIsValid)

	var notes []Note
	err := connection.Db.C("gonotes").Find(nil).All(&notes)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		panic(err)
	}
}
