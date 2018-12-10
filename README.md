# Go API Starter kit

Basic starting point for building an API in Go.
The API routes for notes require a valid JWT token. The token is returned when calling Sign up or Log in.

Note! For improved security use ssh keys to sign the token.

## Setup

### 1. Configuration
Settings for Db server, database and jwt encryption secret are in settings.go. Add this file to .gitingore when using production settings. DO NOT save production keys in Github.

### 2. Install packages and start server
````
$ go get
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

### Todo
- Improve error handling
- Tweak Returns
- Move authentication to middleware
- Status codes as constants
