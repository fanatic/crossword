package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
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
	router.HandleFunc("/layouts", GetLayouts(state)).Methods("GET")
	router.HandleFunc("/layouts", PostLayout(state)).Methods("POST")

	return xray.Handler(
		xray.NewFixedSegmentNamer("crossword-app"),
		handlers.CORS(
			handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"}))(router))
}

// PostGame creates a new game state
func PostGame(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var gameRequest struct {
			BoardID string `json:"board_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&gameRequest); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error parsing body: %s", err)
			return
		}

		id, err := state.CreateGame(r.Context(), gameRequest.BoardID)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error creating game: %s", err)
			return
		}

		game, err := FetchGame(r.Context(), state, id)
		if err != nil {
			if err.Error() == "Timed out" {
				w.WriteHeader(200)
				fmt.Fprintf(w, `{"error": "Timed out"}`)
				return
			}
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
		game, err := FetchGame(r.Context(), state, mux.Vars(r)["id"])
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

		err = state.CreatePlayer(r.Context(), game.ID, playerRequest.Name)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error creating game: %s", err)
			return
		}

		game, err = FetchGame(r.Context(), state, game.ID)
		if err != nil {
			if err.Error() == "Timed out" {
				w.WriteHeader(200)
				fmt.Fprintf(w, `{"error": "Timed out"}`)
				return
			}
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
		game, err := FetchGame(r.Context(), state, mux.Vars(r)["id"])
		if err != nil {
			if err.Error() == "Timed out" {
				w.WriteHeader(200)
				fmt.Fprintf(w, `{"error": "Timed out"}`)
				return
			}
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

		playerExists := false
		for _, player := range game.CurrentPlayers {
			if player.Name == guessRequest.PlayerName {
				playerExists = true
			}
		}
		if !playerExists {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Player does not exist: %s", guessRequest.PlayerName)
			return
		}

		err = state.CreateGuess(r.Context(), game.ID, guessRequest.PlayerName, game.CurrentClue.Number, game.CurrentClue.Direction, guessRequest.Guess)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error creating guess: %s", err)
			return
		}

		// If this player was the last one we were waiting on, continue
		if len(game.CurrentClue.WaitingOnPlayers) == 0 || (len(game.CurrentClue.WaitingOnPlayers) == 1 && game.CurrentClue.WaitingOnPlayers[0] == guessRequest.PlayerName) {
			err := IncrementClue(r.Context(), state, game)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "Error incrementing clue: %s", err)
				return
			}
		}

		game, err = FetchGame(r.Context(), state, game.ID)
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
		game, err := FetchGame(r.Context(), state, mux.Vars(r)["id"])
		if err != nil {
			if err.Error() == "Timed out" {
				w.WriteHeader(200)
				fmt.Fprintf(w, `{"error": "Timed out"}`)
				return
			}
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching game state: %s", err)
			return
		}

		if err := CheckTimers(r.Context(), state, game); err != nil {
			log.Println("Timer error:", err)
		}

		json.NewEncoder(w).Encode(game)
	}
}

// GetLayouts uploads a board
func GetLayouts(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		boardLayouts, err := state.GetBoardLayouts(r.Context())
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching layout: %s", err)
			return
		}

		json.NewEncoder(w).Encode(boardLayouts)
	}
}

// PostLayout uploads a board
func PostLayout(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var layoutRequest struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&layoutRequest); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error parsing body: %s", err)
			return
		}

		rreq, err := http.NewRequest("GET", layoutRequest.URL, nil)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching layout: %s", err)
			return
		}
		rreq.Header.Add("Referer", "https://www.xwordinfo.com/JSON/Sample1")
		req, err := http.DefaultClient.Do(rreq)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching layout: %s", err)
			return
		}
		defer req.Body.Close()

		var boardLayout BoardLayout
		json.NewDecoder(req.Body).Decode(&boardLayout)

		boardLayout.ID = fmt.Sprintf("%s-%s", boardLayout.Publisher, boardLayout.Date)

		err = state.CreateBoardLayout(boardLayout)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error creating layout: %s", err)
			return
		}

		boardLayouts, err := state.GetBoardLayouts(r.Context())
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error fetching layout: %s", err)
			return
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(boardLayouts)
	}
}
