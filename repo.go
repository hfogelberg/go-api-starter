package main

import (
	"log"
)

func (connection *Connection) SaveUser(user User) bool {
	err := connection.Db.C("gousers").Insert(&user)
	if err != nil {
		log.Println("Failed insert")
		return false
	}
	return true
}
