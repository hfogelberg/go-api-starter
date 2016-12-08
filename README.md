#Go API Starter kit

Basic starting point for building an API in Go.
The API routes for notes require a valid JWT token. The token is returned when calling Sign up or Log in.

## Setup
### 1. Install Go packages
````
$ go get github.com/gorilla/mux
$ go get github.com/gorilla/handlers
$ go get github.com/dgrijalva/jwt-go
$ go golang.org/x/crypto/bcrypt
$ go gopkg.in/mgo.v2
$ go gopkg.in/mgo.v2/bson
````
### 2. Configuration
Settings for Db server, database and jwt encryption secret are in config.go. Add this file to .gitingore when using production settings. DO NOT save production keys in Github.

### 3. Start server
````
$ mongod
$ go build
$ go run ./go-api-starter
````

##API calls
To test the API I prefer to use [httpie](https://github.com/jkbrzt/httpie), which I find much easier to use than curl.

### Sign up
````
$ http -f POST http://localhost:8081/api/signup email=batman@mail.com username=Batman password=password
````

### Log in
````
$ http -f POST http://localhost:3000/api/login email=obi@gmail.com username=Obi password=password
````

### Change user

### Remove user

### Create notes
````
$ http -f POST http://localhost:3000/api/noter username=obi text=Lorem Ipsum ...
````

### Get Notes by username
````
$ http GET http://localhost:3000/api/notes
````

### Get note by Id
````
$ http GET http://localhost:3000/api/notes/123
````

### Update note

### Todo
- Status codes as constants
