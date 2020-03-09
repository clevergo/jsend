package main

import (
	"errors"
	"fmt"
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
			{"foo", "foo@example.com"},
			{"bar", "bar@example.com"},
		},
	}
}

type modelError struct {
	errs map[string]error
}

func (me modelError) Error() (s string) {
	for field, err := range me.errs {
		s += fmt.Sprintf("%s: %s\n", field, err)
	}
	return
}

func (me *modelError) Add(field string, err error) {
	if me.errs == nil {
		me.errs = make(map[string]error)
	}

	me.errs[field] = err
}

func (me *modelError) IsEmpty() bool {
	return len(me.errs) == 0
}

func (me *modelError) Errors() map[string]error {
	return me.errs
}

type errorHandler struct {
}

func (eh *errorHandler) Handle(ctx *clevergo.Context, err error) {
	switch e := err.(type) {
	case *modelError:
		errs := e.Errors()
		data := make(map[string]string)
		for field, msg := range errs {
			data[field] = msg.Error()
		}
		err = jsend.Fail(ctx.Response, data)
	default:
		err = jsend.Error(ctx.Response, err.Error(), http.StatusInternalServerError)
	}

	// convert error as jsend response
	if err != nil {
		http.Error(ctx.Response, err.Error(), http.StatusInternalServerError)
	}
}

func getUsers(ctx *clevergo.Context) error {
	return jsend.Success(ctx.Response, users.entries)
}

func getUser(ctx *clevergo.Context) error {
	id := ctx.Params.String("id")
	user, found := users.find(id)
	if !found {
		return errors.New("User Not Found")
	}

	return jsend.Success(ctx.Response, user)
}

func createUser(ctx *clevergo.Context) error {
	if err := ctx.Request.ParseForm(); err != nil {
		return err
	}

	errs := &modelError{}
	id := ctx.Request.FormValue("id")
	if id == "" {
		errs.Add("id", errors.New("id can not be blank"))
	} else if _, found := users.find(id); found {
		errs.Add("id", errors.New("id was taken"))
	}
	email := ctx.Request.FormValue("email")
	if email == "" {
		errs.Add("email", errors.New("email can not be blank"))
	}
	if !errs.IsEmpty() {
		return errs
	}

	user := User{
		ID:    id,
		Email: email,
	}
	users.insert(user)

	return jsend.Success(ctx.Response, user)
}

func deleteUser(ctx *clevergo.Context) error {
	id := ctx.Params.String("id")
	user, found := users.find(id)
	if !found {
		return errors.New("User Not Found")
	}

	users.delete(user.ID)
	return jsend.Success(ctx.Response, nil)
}

func main() {
	router := clevergo.NewRouter()
	router.ErrorHandler = &errorHandler{}
	router.Get("/users", getUsers)
	router.Post("/users", createUser)
	router.Get("/users/:id", getUser)
	router.Delete("/users/:id", deleteUser)
	http.ListenAndServe(":1234", router)
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
