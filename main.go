package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	mgo "gopkg.in/mgo.v2"
)

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
