package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Game struct {
	ID             string      `json:"id" dynamo:"ID,hash"`
	BoardLayout    BoardLayout `json:"-"  dynamo:"-"`
	Grid           []string    `json:"grid" dynamo:"-"`
	GridNums       []int       `json:"grid_nums" dynamo:"-"`
	GridRows       int         `json:"grid_rows" dynamo:"-"`
	GridCols       int         `json:"grid_cols" dynamo:"-"`
	LastClue       *Clue       `json:"last_clue"  dynamo:"-"`
	CurrentClue    Clue        `json:"current_clue"`
	CurrentPlayers []Player    `json:"current_players"`
}

type Player struct {
	Name         string  `json:"name"`
	CurrentScore int     `json:"current_score" dynamo:"-"`
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
	Description string    `json:"description" dynamo:"-"`
	Answer      string    `json:"answer" dynamo:"-"` // ?A?
	Guesses     []Guess   `json:"guesses" dynamo:"-"`
	ExpiresAt   time.Time `json:"expires_at" dynamo:"-"`
}

func GetGame(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		game := Game{
			ID:       "BLAH",
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
		}
		json.NewEncoder(w).Encode(game)
	}
}

func main() {
	state := NewState()

	router := mux.NewRouter()
	router.HandleFunc("/games/{id}", GetGame(state)).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
