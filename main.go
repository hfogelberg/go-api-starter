package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MongoDBHost = "localhost"
	MongoDb     = "test"
	HmacSecret  = "secret"
	Port        = ":3000"
)

type Retval struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Note struct {
	Text     string    `json:"text"`
	Username string    `json:"username"`
	When     time.Time `json:"when" bson:"when"`
}

type Connection struct {
	Db *mgo.Database
}

type User struct {
	Email          string `json:"email"`
	Username       string `json:"username"`
	HashedPassword []byte `json:"password"`
}

type JwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JwtUser struct {
	TokenIsValid bool
	Username     string
}

var hmacSampleSecret = []byte(HmacSecret)

func main() {
	session, err := mgo.Dial(MongoDBHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	connection := Connection{session.DB(MongoDb)}

	router := httprouter.New()

	router.POST("/api/signup", connection.Signup)
	router.POST("/api/login", connection.Login)
	router.GET("/api/notes", connection.GetNotes)
	router.POST("/api/notes", connection.CreateNote)

	log.Fatal(http.ListenAndServe(":3000", router))
}

// Helpers
func tokenIsValid(tokenString string) bool {
	isValid := true

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token validated")
		fmt.Println(claims)
		fmt.Println(claims["username"])
	} else {
		fmt.Println(err)
		isValid = false
	}

	return isValid
}

// Handlers
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

	log.Println("Signup")

	// Parse body and hash password
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

	log.Println("User", user)

	err = connection.Db.C("gousers").Insert(&user)
	if err != nil {
		log.Println("Failed insert")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := CreateToken(user.Username)
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

func CreateToken(username string) string {
	log.Println("CreateToken")
	expireToken := time.Now().Add(time.Minute * 60).Unix()

	claims := JwtClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:3000",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var hmacSampleSecret = []byte(HmacSecret)
	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		log.Println("Error signing token ", err)
	}

	log.Println("Token created ", tokenString)

	return tokenString
}
