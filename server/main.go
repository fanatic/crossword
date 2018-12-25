package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Game struct {
	ID string `json:"id"`
}

func GetGame(w http.ResponseWriter, r *http.Request) {
	game := Game{ID: "BLAH"}
	json.NewEncoder(w).Encode(game)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/games/{id}", GetGame).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
