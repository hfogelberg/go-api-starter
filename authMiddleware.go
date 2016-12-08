package main

import (
	"log"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var hmacSampleSecret = []byte(HmacSecret)

// func validateToken(tokenString string) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return hmacSampleSecret, nil
// 	})
//
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		fmt.Println("Token validated")
// 		fmt.Println(claims)
// 		fmt.Println(claims["username"])
// 	} else {
// 		fmt.Println(err)
// 	}
// }

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		log.Println("Middleware")
		return HmacSecret, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// *******************************
// Koden jag fick kontakt med !!!!!!!!
// func ValidateToken(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("middleware", r.URL)
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return hmacSampleSecret, nil
// 		})
//
// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			fmt.Println("Token validated")
// 			fmt.Println(claims)
// 			fmt.Println(claims["username"])
// 		} else {
// 			fmt.Println(err)
// 		}
//
// 		h.ServeHTTP(w, r)
// 	})
// }
// ***********************************

// func Middleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("middleware", r.URL)
// 		h.ServeHTTP(w, r)
// 	})
// }
