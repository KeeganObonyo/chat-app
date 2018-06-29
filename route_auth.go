package main

import (
	"./data"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

var errormessage string = "POST"
var successmessage string = "User created successfully"

// POST /signup
//Create a new user
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	user := data.User{}
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Data limit exceeded")
	}
	if err := request.Body.Close(); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Invalid request body")
	}
	if err := json.Unmarshal(body, &user); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(writer).Encode(errormessage); err != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(400)
			json.NewEncoder(writer).Encode("Invalid json data")
		}
	}
	if err := user.Create(); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Couldn't create user")
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		json.NewEncoder(writer).Encode(successmessage)
	}
}

// POST /authenticate
// Authenticate the user given the email and password

type Login struct {
	email    string
	password string
}

func (f Login) Email() string {
	return f.email
}
func (f Login) Password() string {
	return f.password
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {
	login := Login{}
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Data limit exceeded")
	}
	if err := request.Body.Close(); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Invalid request body")
	}
	if err := json.Unmarshal(body, &login); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(writer).Encode(errormessage); err != nil {
			danger(err, "errors")
		}
	}
	// err := request.ParseForm()
	user, err := data.UserByEmail(login.Email())
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Couldn't find user")
	}
	if user.Password == data.Encrypt(login.Email()) {
		session, err := user.CreateSession()
		if err != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(400)
			json.NewEncoder(writer).Encode("Error creating session")
		} else {
			cookie := http.Cookie{
				Name:     "_cookie",
				Value:    session.Uuid,
				HttpOnly: true,
			}
			http.SetCookie(writer, &cookie)
		}
	}

}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
}
