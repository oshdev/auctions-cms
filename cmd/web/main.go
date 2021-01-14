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
	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal("cannot listen and serve", err)
	}
}
