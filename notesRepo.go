package main

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

func (connection *Connection) SaveNoteInDb(note Note) bool {
	// TODO Fix error codes

	if err := connection.Db.C("gonotes").Insert(&note); err != nil {
		return false
	}

	return true
}

func (connection *Connection) GetNotesByUser(username string) []Note {
	//TODO error handling
	log.Println("Get notes for user ", username)

	var notes []Note
	if err := connection.Db.C("gonotes").Find(bson.M{"username": username}).All(&notes); err != nil {
		log.Println("Error fetching notes for user with username ", username)
		log.Println(err)
	}

	return notes
}

func (connection *Connection) GetNoteById(noteId string) Note {
	//TODO error handling

	var note Note
	if err := connection.Db.C("gonotes").Find(bson.M{"_id": noteId}).One(&note); err != nil {
		log.Println("Error fetching note with id ", noteId)
		log.Println(err)
	}

	return note
}

func (connection *Connection) RemoveNote(noteId string) bool {
	if err := connection.Db.C("gonotes").RemoveId(noteId); err != nil {
		log.Println("Error deleting document with id ", noteId)
		log.Println(err)
		return false
	}

	return true
}
