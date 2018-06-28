package main

import (
	// "./data"
	"encoding/json"
	// "fmt"
	"net/http"
)

const (
	empty = ""
	tab   = "\t"
)

var message1 string = "Internal functionality error"

func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(vals)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(vals)
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	urls := map[int]string{
		0: "Welcome",
		1: "/logout",
		2: "/signup",
		3: "/authenticate",
		4: "/thread/create",
		5: "/thread/post",
		6: "/thread/read",
	}
	{
		writer.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(urls)

	}

}
