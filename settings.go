package main

const (
	MongoDBHost = "localhost"
	MongoDb     = "test"
	HmacSecret  = "secret"
	Port        = ":3000"
)

var hmacSampleSecret = []byte(HmacSecret)
