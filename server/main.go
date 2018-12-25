package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Game struct {
	ID             string      `json:"id"`
	BoardLayout    BoardLayout `json:"-"`
	Grid           []string    `json:"grid"`
	GridNums       []int       `json:"grid_nums"`
	GridRows       int         `json:"grid_rows"`
	GridCols       int         `json:"grid_cols"`
	LastClue       Clue        `json:"last_clue"`
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
	Clue        Clue      `json:"-" dynamo:"-"`
	Guess       string    `json:"guess"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type Clue struct {
	Number           int        `json:"number"`
	Direction        string     `json:"direction"`
	Description      string     `json:"description"`
	Answer           string     `json:"answer"` // ?A?
	Guesses          []Guess    `json:"guesses"`
	WaitingOnPlayers []string   `json:"waiting_on_players"`
	ExpiresAt        *time.Time `json:"expires_at"`
}

func WaitingOnPlayers(currentPlayers []Player, guesses []Guess) []string {
	players := map[string]bool{}
	for _, p := range currentPlayers {
		players[p.Name] = true
	}
	for _, g := range guesses {
		delete(players, g.Player.Name)
	}

	waitingOnPlayers := []string{}
	for player := range players {
		waitingOnPlayers = append(waitingOnPlayers, player)
	}
	return waitingOnPlayers
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

	guesses, err := state.GetGuesses(gameID)
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

	clueGuesses := map[string][]Guess{}
	playerGuesses := map[string][]Guess{}
	for _, guess := range guesses {
		g := Guess{
			Clue:        Clue{Number: guess.ClueNumber(), Direction: guess.ClueDirection()},
			Player:      Player{Name: guess.PlayerName()},
			Guess:       guess.Guess,
			SubmittedAt: guess.UpdatedAt,
		}
		clueID := fmt.Sprintf("%d-%s", guess.ClueNumber(), guess.ClueDirection())
		clueGuesses[clueID] = append(clueGuesses[clueID], g)
		playerGuesses[guess.PlayerName()] = append(playerGuesses[guess.PlayerName()], g)
	}

	g := Game{
		ID:             game.ID,
		Grid:           grid,
		GridNums:       sampleBoard.Gridnums,
		GridRows:       sampleBoard.Size.Rows,
		GridCols:       sampleBoard.Size.Cols,
		CurrentClue:    sampleBoard.GetClue(game.CurrentClueNumber, game.CurrentClueDirection, &game.CurrentClueExpiresAt),
		LastClue:       sampleBoard.GetLastClue(game.CurrentClueNumber, game.CurrentClueDirection),
		CurrentPlayers: []Player{},
	}
	lastClueID := fmt.Sprintf("%d-%s", g.LastClue.Number, g.LastClue.Direction)
	g.LastClue.Guesses = clueGuesses[lastClueID]

	for _, player := range players {
		g.CurrentPlayers = append(g.CurrentPlayers, Player{
			Name:         player.PlayerName,
			CurrentScore: 0,
			Guesses:      playerGuesses[player.PlayerName],
		})
	}

	currentClueID := fmt.Sprintf("%d-%s", g.CurrentClue.Number, g.CurrentClue.Direction)
	g.CurrentClue.WaitingOnPlayers = WaitingOnPlayers(g.CurrentPlayers, clueGuesses[currentClueID])

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

func IncrementClue(state *State, game *Game) error {
	number, direction := sampleBoard.Next(game.CurrentClue.Number, game.CurrentClue.Direction)
	return state.UpdateGameClue(game.ID, number, direction)
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
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"}))(router)))
}
