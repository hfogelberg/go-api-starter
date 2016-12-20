package main

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func tokenIsValid(tokenString string) bool {
	isValid := true

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
