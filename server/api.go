package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	state := NewState()

	router := mux.NewRouter()
	router.HandleFunc("/games", PostGame(state)).Methods("POST")
	router.HandleFunc("/games/{id}", GetGame(state)).Methods("GET")
	router.HandleFunc("/games/{id}/players", PostPlayer(state)).Methods("POST")
	router.HandleFunc("/games/{id}/guesses", PostGuess(state)).Methods("POST")
	return handlers.CORS(handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"}))(router)
}

// PostGame creates a new game state
func PostGame(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := state.CreateGame()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error creating game: %s", err)
			return
		}

		game, err := FetchGame(state, id)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching game state: %s", err)
			return
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(game)
	}
}

// PostPlayer creates new player in game
func PostPlayer(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		game, err := FetchGame(state, mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching game state: %s", err)
			return
		}

		var playerRequest struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error parsing body: %s", err)
			return
		}

		err = state.CreatePlayer(game.ID, playerRequest.Name)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error creating game: %s", err)
			return
		}

		game, err = FetchGame(state, game.ID)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching game state: %s", err)
			return
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(game)
	}
}

// PostGuess creates new guess
func PostGuess(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		game, err := FetchGame(state, mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching game state: %s", err)
			return
		}

		var guessRequest struct {
			PlayerName string `json:"player_name"`
			Guess      string `json:"guess"`
		}
		if err := json.NewDecoder(r.Body).Decode(&guessRequest); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error parsing body: %s", err)
			return
		}

		err = state.CreateGuess(game.ID, guessRequest.PlayerName, game.CurrentClue.Number, game.CurrentClue.Direction, guessRequest.Guess)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error creating guess: %s", err)
			return
		}

		if len(game.CurrentClue.WaitingOnPlayers) == 1 && game.CurrentClue.WaitingOnPlayers[0] == guessRequest.PlayerName {
			err := IncrementClue(state, game)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "Error incrementing clue: %s", err)
				return
			}
		}

		game, err = FetchGame(state, game.ID)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching game state: %s", err)
			return
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(game)
	}
}

// GetGame gets current game state
func GetGame(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		game, err := FetchGame(state, mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching game state: %s", err)
			return
		}

		json.NewEncoder(w).Encode(game)
	}
}
