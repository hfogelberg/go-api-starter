package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"golang.org/x/crypto/bcrypt"
)

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		signup(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	var userInDb User
	var user User

	// Parse body and hash password
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	password := r.Form["password"][0]

	pwd := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pwd, 10)
	log.Println("Hashed password", reflect.TypeOf(hashedPassword))

	user.Email = r.Form["email"][0]
	user.Username = r.Form["username"][0]
	user.HashedPassword = hashedPassword

	// Hook up to Db
	db := context.Get(r, "database").(*mgo.Session)

	// Check if username is already in use
	err = db.DB("test").C("gousers").Find(bson.M{"username": user.Username}).One(&userInDb)
	if err == nil {
		log.Println("Username already taken")

		ret := Retval{
			Status:  -100,
			Message: "Username already taken",
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
		return
	}

	// Check if email is already in Db
	log.Println("Now time to check email")
	err = db.DB("test").C("gousers").Find(bson.M{"email": user.Email}).One(&userInDb)
	if err == nil {
		log.Println("Email is already in Db")
		ret := Retval{
			Status:  -101,
			Message: "You have already signed up with that email address",
		}

		// There's aleady a user wth the email address
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
		return
	}

	// insert new user into the database and return token
	if err := db.DB("test").C("gousers").Insert(&user); err != nil {
		log.Println("Failed insert")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Insert OK")

	token := createToken(user.Username)
	log.Println("We have a token!", token)

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

}
