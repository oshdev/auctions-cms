package main

import (
	todo "go-hotwire"
	"log"
	"net/http"
)

func main() {
	server, err := todo.NewServer("../../html/*", todo.NewInMemoryRepo())
	if err != nil {
		log.Fatal(err)
	}
	http.ListenAndServe(":8080", server)
}
