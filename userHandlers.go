package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (connection *Connection) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]

	user := connection.UsernameIsInDb(username)
	if user.Username != "" {
		pwd := []byte(password)
		err = bcrypt.CompareHashAndPassword(user.HashedPassword, pwd)
		fmt.Println(err) // nil means it is a match
		if err == nil {
			// Password OK. Generate token
			token := CreateToken(username)

			ret := Retval{
				Status:  100,
				Token:   token,
				Message: "OK",
			}

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

func (connection *Connection) Signup(w http.ResponseWriter, r *http.Request) {
	var user User

	// TODO Check that username and email is not in DB already

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
