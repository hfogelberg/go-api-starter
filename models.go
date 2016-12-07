package main

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Email          string `json:"email"`
	Username       string `json:"username"`
	HashedPassword []byte `json:"password"`
}

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

type JwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Notes []Note
