package main

import (
	"./data"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Thread struct {
	topic string
}

// SetName receives a pointer to Thread so it can modify it.
func (f *Thread) SetName(topic string) {
	f.topic = topic
}

// Name receives a copy of Thread since it doesn't need to modify it.
func (f Thread) Name() string {
	return f.topic
}

func createThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	thread := Thread{}
	if err != nil {
		error_message(writer, request, "Invalid session")
	} else {

		body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
		if err != nil {
			error_message(writer, request, "Data limit exceeded")
		}
		if err := request.Body.Close(); err != nil {
			error_message(writer, request, "Invalid data")
		}
		if err := json.Unmarshal(body, &thread); err != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(422)
			json.NewEncoder(writer).Encode("Invalid json data")
		} else {
			topic := thread.Name()
			user, err := sess.User()
			if err != nil {
				danger(err, "Cannot get user from session")
			}
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)
			if _, err := user.CreateThread(topic); err != nil {
				danger(err, "Cannot create thread")
			}
		}
	}
}

// GET /thread/read
func readThread(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	thread, err := data.ThreadByUUID(uuid)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Cannot read thread")
	} else {
		_, err := session(writer, request)
		if err != nil {
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode("Invalid session")
		} else {
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(thread)
		}
	}
}

// POST /thread/post
// Create the post

type PostThread struct {
	body string
	uuid string
}

func (f PostThread) Uuid() string {
	return f.uuid
}

func (f PostThread) Body() string {
	return f.body
}
func postThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	postthread := PostThread{}
	if err != nil {
		error_message(writer, request, "Invalid session")
	} else {

		body, err := ioutil.ReadAll(io.LimitReader(request.Body, 104857655))
		if err != nil {
			error_message(writer, request, "Data limit excceeded")
		}
		if err := request.Body.Close(); err != nil {
			error_message(writer, request, "Invalid request body")
		}
		if err := json.Unmarshal(body, &postthread); err != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(422)
			json.NewEncoder(writer).Encode("Invalid Json")
		} else {
			user, err := sess.User()
			if err != nil {
				danger(err, "Cannot get user from session")
			}

			postbody := postthread.Body()
			uuid := postthread.Uuid()

			thread, err := data.ThreadByUUID(uuid)
			if err != nil {
				error_message(writer, request, "No previous thread to embedd")
			}
			if _, err := user.CreatePost(thread, postbody); err != nil {
				danger(err, "Cannot create post")
			}
		}
	}
}
