package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
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

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/signup", connection.Signup)
	router.HandleFunc("/api/login", connection.Login)

	router.HandleFunc("/api/notes", connection.GetNotes)
	router.HandleFunc("/api/createnote", connection.CreateNote)

	// n := negroni.Classic()
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)

	n.UseHandler(router)
	n.Run(":3000")

	// log.Fatal(http.ListenAndServe(":3000", router))
}
