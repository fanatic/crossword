package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Game struct {
	ID             string      `json:"id"`
	BoardLayout    BoardLayout `json:"-"`
	Grid           []string    `json:"grid"`
	GridNums       []int       `json:"grid_nums"`
	GridRows       int         `json:"grid_rows"`
	GridCols       int         `json:"grid_cols"`
	LastClue       *Clue       `json:"last_clue"`
	CurrentClue    Clue        `json:"current_clue"`
	CurrentPlayers []Player    `json:"current_players"`
}

type Player struct {
	Name         string  `json:"name"`
	CurrentScore int     `json:"current_score"`
	Guesses      []Guess `json:"-"`
}

type Guess struct {
	Player      Player    `json:"player" dynamo:"-"`
	Clue        Clue      `json:"clue" dynamo:"-"`
	Guess       string    `json:"guess"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type Clue struct {
	Number      int       `json:"number"`
	Direction   string    `json:"direction"`
	Description string    `json:"description"`
	Answer      string    `json:"answer"` // ?A?
	Guesses     []Guess   `json:"guesses"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func FetchGame(state *State, gameID string) (*Game, error) {
	game, err := state.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	players, err := state.GetPlayers(gameID)
	if err != nil {
		return nil, err
	}

	grid := make([]string, len(sampleBoard.Grid))
	for i := range sampleBoard.Grid {
		switch sampleBoard.Grid[i] {
		case ".":
			grid[i] = "."
			break
		default:
			grid[i] = " "
		}
	}
	g := Game{
		ID:       game.ID,
		Grid:     grid,
		GridNums: sampleBoard.Gridnums,
		GridRows: sampleBoard.Size.Rows,
		GridCols: sampleBoard.Size.Cols,
		CurrentClue: Clue{
			Number:      1,
			Direction:   "across",
			Description: sampleBoard.Clues.Across[0],
			Answer:      sampleBoard.Answers.Across[0],
			ExpiresAt:   time.Now().Add(30 * time.Second),
		},
		CurrentPlayers: []Player{},
	}

	for _, player := range players {
		guesses := []Guess{}
		for clueID, guess := range player.Guesses {
			clueNum, _ := strconv.Atoi(strings.Split(clueID, "-")[0])
			clueDir := strings.Split(clueID, "-")[1]
			guesses = append(guesses, Guess{
				Clue: Clue{
					Number:    clueNum,
					Direction: clueDir,
				},
				Guess:       guess.Guess,
				SubmittedAt: guess.SubmittedAt,
			})
		}

		g.CurrentPlayers = append(g.CurrentPlayers, Player{
			Name:         player.PlayerName,
			CurrentScore: 0,
			Guesses:      guesses,
		})
	}

	return &g, nil
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

func main() {
	state := NewState()

	router := mux.NewRouter()
	router.HandleFunc("/games", PostGame(state)).Methods("POST")
	router.HandleFunc("/games/{id}", GetGame(state)).Methods("GET")
	router.HandleFunc("/games/{id}/players", PostPlayer(state)).Methods("POST")
	router.HandleFunc("/games/{id}/guesses", PostGuess(state)).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
