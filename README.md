# JSend's implementation writen in Go(aka golang)
[![Build Status](https://travis-ci.org/clevergo/jsend.svg?branch=master)](https://travis-ci.org/clevergo/jsend) [![Coverage Status](https://coveralls.io/repos/github/clevergo/jsend/badge.svg?branch=master)](https://coveralls.io/github/clevergo/jsend?branch=master) [![GoDoc](https://godoc.org/github.com/clevergo/jsend?status.svg)](http://godoc.org/github.com/clevergo/jsend) [![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/jsend)](https://goreportcard.com/report/github.com/clevergo/jsend) [![Release](https://img.shields.io/github/release/clevergo/jsend.svg?style=flat-square)](https://github.com/clevergo/jsend/releases)

This package is an implementation of [JSend](https://github.com/omniti-labs/jsend) written in Go(aka golang).

## Installation

```shell
go get github.com/clevergo/jsend
```

## Usage

Usage is pretty simple.

```go
// success response
jsend.Success(w, "any type of data")
// fail response
jsend.Fail(w, "any type of data")
// error response
jsend.Error(w, "error message")
// error response with extra code
jsend.ErrorCode(w, "error message", errorCode)
// error response with extra code and data
jsend.ErrorCodeData(w, "error message", errorCode, "any type of data")
```

### Error Handling

It is application responsibility to handle error, see [Full Example](#full-example) for further more detail.

### Status Code

By default status code `http.StatusOK` was used implicitly,
it can also be specified by the last parameter if necessary.

```go
jsend.Success(w, "any type of data", http.StatusOK)
jsend.Fail(w, "any type of data", http.StatusForbidden)
jsend.Error(w, "error message", http.StatusInternalServerError)
```

### Full Example

```go
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/jsend"
)

var users *Users

func init() {
	users = &Users{
		mutex: &sync.RWMutex{},
		entries: []User{
			User{"foo", "foo@example.com"},
			User{"bar", "bar@example.com"},
		},
	}
}

func handleError(w http.ResponseWriter, err error) {
	// for method chaining
	if err == nil {
		return
	}

	log.Println(err.Error())
	// convert error as jsend response
	if err = jsend.Error(w, err.Error(), http.StatusInternalServerError); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	handleError(w, jsend.Success(w, users.entries))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := clevergo.GetParams(r).Get("id")
	user, found := users.find(id)
	if !found {
		handleError(w, jsend.Error(w, "User Not Found", http.StatusNotFound))
		return
	}

	handleError(w, jsend.Success(w, user))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, err)
		return
	}

	errs := map[string][]string{}
	id := r.FormValue("id")
	if id == "" {
		errs["id"] = append(errs["id"], "id can not be blank")
	}
	if _, found := users.find(id); found {
		errs["id"] = append(errs["id"], "id was taken")
	}
	email := r.FormValue("email")
	if email == "" {
		errs["email"] = append(errs["email"], "email can not be blank")
	}
	if len(errs) > 0 {
		handleError(w, jsend.Fail(w, errs))
		return
	}

	user := User{
		ID:    id,
		Email: email,
	}
	users.insert(user)

	handleError(w, jsend.Success(w, user))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := clevergo.GetParams(r).Get("id")
	user, found := users.find(id)
	if !found {
		handleError(w, jsend.Error(w, "User Not Found", http.StatusNotFound))
		return
	}

	users.delete(user.ID)
	handleError(w, jsend.Success(w, nil))
}

func main() {
	app := clevergo.New(":1234")
	app.Get("/users", getUsers)
	app.Post("/users", createUser)
	app.Get("/users/:id", getUser)
	app.Delete("/users/:id", deleteUser)
	app.ListenAndServe()
}

type Users struct {
	entries []User
	mutex   *sync.RWMutex
}

func (us *Users) find(id string) (User, bool) {
	for _, user := range users.entries {
		if user.ID == id {
			return user, true
		}
	}

	return User{}, false
}

func (us *Users) insert(user User) {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	us.entries = append(us.entries, user)
}

func (us *Users) delete(id string) {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	for i, user := range us.entries {
		if user.ID == id {
			us.entries = append(us.entries[:i], us.entries[i+1:]...)
		}
	}
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
```

```shell
$ go run main.go

# fetch users
$ curl http://localhost:1234/users
{"status":"success","data":[{"id":"foo","email":"foo@example.com"},{"id":"bar","email":"bar@example.com"}]}

# create user without required data
$ curl -XPOST  http://localhost:1234/users
{"status":"fail","data":{"email":["email can not be blank"],"id":["id can not be blank"]}}

# create test user
$ curl -XPOST -d "id=test&email=test@example.com" http://localhost:1234/users
{"status":"success","data":{"id":"test","email":"test@example.com"}}

# refetch users
$ curl http://localhost:1234/users
{"status":"success","data":[{"id":"foo","email":"foo@example.com"},{"id":"bar","email":"bar@example.com"},{"id":"test","email":"test@example.com"}]}

# fetch test user
$ curl http://localhost:1234/users/test
{"status":"success","data":{"id":"test","email":"test@example.com"}}

# delete test user
$ curl -XDELETE http://localhost:1234/users/test
{"status":"success","data":null}

# refetch test user
$ curl http://localhost:1234/users/test
{"status":"error","data":null,"message":"User Not Found"}
```
