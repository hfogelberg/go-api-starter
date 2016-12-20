#Go API Starter kit

Basic starting point for building an API in Go.
The API routes for notes require a valid JWT token. The token is returned when calling Sign up or Log in.

## Setup
### 1. Install Go packages
````
$ go get github.com/julienschmidt/httprouter
$ go get github.com/dgrijalva/jwt-go
$ go golang.org/x/crypto/bcrypt
$ go gopkg.in/mgo.v2
$ go gopkg.in/mgo.v2/bson
````
### 2. Configuration
Settings for Db server, database and jwt encryption secret are in settings.go. Add this file to .gitingore when using production settings. DO NOT save production keys in Github.

### 3. Start server
````
$ mongod
$ go build
$ go run ./go-api-starter
````


### Sign up
````
POST http://localhost:8081/api/signup
Body: email, username, password
Returns token
````

### Log in
````
POST http://localhost:3000/api/login
Body: username, password
Returns token
````

### Change user
Todo

### Remove user
Todo

### Create notes
````
POST: http://localhost:3000/api/notes
Authentication: token
Body: Text
````

### Get Notes by username
````
GET http://localhost:3000/api/notes
Authentication: token
````

### Get note by Id
````
GET http://localhost:3000/api/notes/123
Authentication: token
````

### Update note
Todo

### Todo
- Status codes as constants
