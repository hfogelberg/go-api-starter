package main

import (
	"log"
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
)

type Adapter func(http.Handler) http.Handler

func main() {
	// connect to the database
	db, err := mgo.Dial(MongoDBHost)
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	// clean up when done
	defer db.Close()

	// Create handlers
	userHandler := Adapt(http.HandlerFunc(handleUser), withDB(db))
	loginHandler := Adapt(http.HandlerFunc(handleLogin), withDB(db))
	notesHandler := Adapt(http.HandlerFunc(handleNotes), withDB(db))

	// add the handler
	http.Handle("/user", context.ClearHandler(userHandler))
	http.Handle("/login", context.ClearHandler(loginHandler))
	http.Handle("/notes", context.ClearHandler(notesHandler))

	// start the server
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func withDB(db *mgo.Session) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbsession := db.Copy()
			defer dbsession.Close()
			context.Set(r, "database", dbsession)
			h.ServeHTTP(w, r)
		})
	}
}

// package main
//
// import (
// 	"log"
// 	"net/http"
//
// 	"github.com/gorilla/context"
// 	mgo "gopkg.in/mgo.v2"
// )
//
// type Adapter func(http.Handler) http.Handler
//
// func main() {
// 	// connect to the database
// 	db, err := mgo.Dial(MongoDBHost)
// 	if err != nil {
// 		log.Println("***** Cannot connect to Mongo *****")
// 		log.Fatal("cannot dial mongo", err)
// 	}
//
// 	// clean up when done
// 	defer db.Close()
//
// 	// Create handlers
// 	userHandler := Adapt(http.HandlerFunc(handleUser), withDB(db))
// 	loginHandler := Adapt(http.HandlerFunc(handleLogin), withDB(db))
// 	notesHandler := Adapt(http.HandlerFunc(handleNotes), withDB(db))
//
// 	// add the handler
// 	http.Handle("/api/user", context.ClearHandler(userHandler))
// 	http.Handle("/api/login", context.ClearHandler(loginHandler))
// 	http.Handle("/api/notes", jwtMiddleware.Handler(context.ClearHandler(notesHandler)))
//
// 	// http.Handle("/notes", Middleware(context.ClearHandler(notesHandler)))
//
// 	// start the server
// 	if err := http.ListenAndServe(Port, nil); err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
// 	for _, adapter := range adapters {
// 		h = adapter(h)
// 	}
// 	return h
// }
//
// func withDB(db *mgo.Session) Adapter {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			dbsession := db.Copy()
// 			defer dbsession.Close()
// 			context.Set(r, "database", dbsession)
// 			h.ServeHTTP(w, r)
// 		})
// 	}
// }
