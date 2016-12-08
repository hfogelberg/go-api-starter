package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		login(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]

	db := context.Get(r, "database").(*mgo.Session)

	// Check if username is already in use
	user := User{}
	err = db.DB(MongoDb).C("gousers").Find(bson.M{"username": username}).One(&user)

	log.Println("User is in Db", user.HashedPassword)
	if user.Username != "" {
		// Compare to password in Db
		pwd := []byte(password)
		err = bcrypt.CompareHashAndPassword(user.HashedPassword, pwd)
		fmt.Println(err) // nil means it is a match
		if err == nil {
			// Password OK. Generate token
			log.Println("Password is OK, Time to generate token")
			token := createToken(username)
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
