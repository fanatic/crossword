package main

import (
	"github.com/akrylysov/algnhsa"
	"github.com/fanatic/crossword/server"
)

func main() {
	algnhsa.ListenAndServe(server.NewRouter(), nil)
}
