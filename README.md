#Go API Starter kit

Basic starting point for building an API in Go.

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
Todo

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
http -f POST http://localhost:8081/api/signup email=batman@mail.com username=Batman password=password
````

### Log in
$ http -f POST http://localhost:8081/api/login email=obi@gmail.com username=Obi password=password

### Change username


### Remove user

### Get Notes by username

### Get notes by Id

### Update note

### Todo
- Status codes as constants
- Settings in seperate file
