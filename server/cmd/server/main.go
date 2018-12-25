package main

import (
	"log"
	"net/http"

	"github.com/fanatic/crossword/server"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", server.NewRouter()))
}
