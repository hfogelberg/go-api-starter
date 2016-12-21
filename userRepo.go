package main

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

func (connection *Connection) SaveUser(user User) bool {
	// TODO Error handling
	err := connection.Db.C("gousers").Insert(&user)
	if err != nil {
		log.Println("Failed insert")
		return false
	}
	return true
}

func (connection *Connection) UsernameIsInDb(username string) User {
	// TODO Error handling
	user := User{}
	err := connection.Db.C("gousers").Find(bson.M{"username": username}).One(&user)

	if err != nil {
		log.Println(err)
	}

	return user
}

func (connection *Connection) EmailIsInDb(email string) User {
	// TODO Error handling
	user := User{}
	err := connection.Db.C("gousers").Find(bson.M{"email": email}).One(&user)
	if err != nil {
		log.Println(err)
	}

	return user
}
