package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Note struct {
	Title string `json:"title"'`
	Text  string `json:"text"'`
}

type JwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type User struct {
	Email          string
	Username       string
	HashedPassword []byte
}

type Retval struct {
	Status  int    `json:status`
	Message string `json:message`
	Token   string `json:token`
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

	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer session.Close()

	// Check if user already exists
	connection := session.DB("test").C("gousers")
	user := User{}
	err = connection.Find(bson.M{"username": username}).One(&user)
	if user.Email == "" {
		// Add new user to Db
		// Encrypt password
		pwd := []byte(password)
		hashedPassword, err := bcrypt.GenerateFromPassword(pwd, 10)
		log.Println("Hashed password", reflect.TypeOf(hashedPassword))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("The hashed password is : %s\n", string(hashedPassword))
		log.Println("Type of Hashed password", reflect.TypeOf(hashedPassword))
		log.Println("Before insert")
		connection := session.DB("test").C("gousers")
		err2 := connection.Insert(&User{email, username, hashedPassword})
		if err2 != nil {
			log.Println("Error saveing to Db", err2)
		}
		log.Println("Insert OK")

		token := createToken(username)
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

	} else {
		ret := Retval{
			Status:  -100,
			Message: "Unauthorized",
		}

		// There's aleady a user wth the email address
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
	}
}

func createToken(username string) string {
	expireToken := time.Now().Add(time.Minute * 1).Unix()

	claims := JwtClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:3000",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		panic(err)
	}

	return tokenString
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]

	// Check if user exists in Db
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer session.Close()

	connection := session.DB("test").C("gousers")
	user := User{}
	err = connection.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Passord in Db", user.HashedPassword)
	if user.Email != "" {
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
