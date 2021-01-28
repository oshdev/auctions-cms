package main

import (
	auction "go-hotwire"
	"log"
	"net/http"
)

const addr = ":8080"

func main() {
	server, err := auction.NewServer("../../html/*", auction.NewInMemoryRepo())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal("cannot listen and serve", err)
	}
}
