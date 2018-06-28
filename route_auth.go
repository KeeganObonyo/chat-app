package main

import (
	"./data"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// GET /login
// Show the login page

// POST /signup
// Create the user account
type NewUser struct {
	name     string
	email    string
	password string
}

func (f NewUser) Name() string {
	return f.name
}
func (f NewUser) Email() string {
	return f.email
}
func (f NewUser) Password() string {
	return f.password
}
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	newuser := NewUser{}
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	if err != nil {
		danger(err, "No data posted")
	}
	if err := request.Body.Close(); err != nil {
		danger(err, "errors")
	}
	if err := json.Unmarshal(body, &newuser); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(writer).Encode(err); err != nil {
			danger(err, "errors")
		}
	}
	user := data.User{
		Name:     newuser.Name(),
		Email:    newuser.Email(),
		Password: newuser.Password(),
	}
	t := user.Create()
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(t); err != nil {
		danger(err, "Couldn't create User")
	}
	http.Redirect(writer, request, "/authenticate", 302)
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		danger(err, "Cannot find user")
	}
	if user.Password == data.Encrypt(request.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, request, "/", 302)
	} else {
		http.Redirect(writer, request, "/signup", 302)
	}

}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		warning(err, "Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}
