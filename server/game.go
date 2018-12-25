package server

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	CurrentScore int     `json:"current_score,omitempty"`
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
	Answer           string     `json:"answer,omitempty"` // ?A?
	Guesses          []Guess    `json:"guesses,omitempty"`
	WaitingOnPlayers []string   `json:"waiting_on_players,omitempty"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
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

	grid := InitializeGrid(sampleBoard.Size.Rows, sampleBoard.Size.Cols, sampleBoard.Grid, sampleBoard.Gridnums)
	for clueID, guesses := range clueGuesses {
		n, _ := strconv.Atoi(strings.Split(clueID, "-")[0])
		d := strings.Split(clueID, "-")[1]
		for _, guess := range guesses {
			// If guess is correct, fill in grid with answer
			answer := strings.ToUpper(guess.Guess)
			if (d == "across" && answer == sampleBoard.Answers.Across[n]) || (d == "down" && answer == sampleBoard.Answers.Down[n]) {
				clueLabel := sampleBoard.ClueLabel(n, d)
				row, col := sampleBoard.ClueLabelPosition(clueLabel)
				if d == "across" {
					grid = FillInGridAnswerAcross(row, col, sampleBoard.Size.Rows, sampleBoard.Size.Cols, answer, grid)
				} else {
					grid = FillInGridAnswerDown(row, col, sampleBoard.Size.Rows, sampleBoard.Size.Cols, answer, grid)
				}
			}
		}
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
	g.LastClue.Answer = ""

	for _, player := range players {
		g.CurrentPlayers = append(g.CurrentPlayers, Player{
			Name:         player.PlayerName,
			CurrentScore: 0,
			Guesses:      playerGuesses[player.PlayerName],
		})
	}

	currentClueID := fmt.Sprintf("%d-%s", g.CurrentClue.Number, g.CurrentClue.Direction)
	g.CurrentClue.WaitingOnPlayers = WaitingOnPlayers(g.CurrentPlayers, clueGuesses[currentClueID])
	g.CurrentClue.Answer = MaskAnswer(g.CurrentClue.Answer, g.CurrentClue.Number, g.CurrentClue.Direction, grid)

	return &g, nil
}

func MaskAnswer(answer string, currentClueNumber int, currentClueDirection string, grid []string) string {
	clueLabel := sampleBoard.ClueLabel(currentClueNumber, currentClueDirection)
	row, col := sampleBoard.ClueLabelPosition(clueLabel)
	numRows, numCols := sampleBoard.Size.Rows, sampleBoard.Size.Cols

	a := ""
	for i := range answer {
		idx := row*numRows + (col+i)%numCols
		if currentClueDirection == "down" {
			idx = (row+i)*numRows + col%numCols
		}
		if grid[idx] == " " {
			a += "?"
		} else {
			a += grid[idx]
		}
	}
	return a
}

func InitializeGrid(numRows, numCols int, grid []string, gridNums []int) []string {
	g := make([]string, len(grid))
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			idx := i*numRows + j%numCols
			switch grid[idx] {
			case ".":
				g[idx] = "."
				break
			default:
				g[idx] = " "
			}
		}
	}
	return g
}

func FillInGridAnswerAcross(startRow, startCol, numRows, numCols int, answer string, grid []string) []string {
	for i := 0; i < len(answer); i++ {
		row := startRow
		col := startCol + i
		grid[row*numRows+col%numCols] = answer[i : i+1]
	}
	return grid
}

func FillInGridAnswerDown(startRow, startCol, numRows, numCols int, answer string, grid []string) []string {
	for i := 0; i < len(answer); i++ {
		row := startRow + i
		col := startCol
		grid[row*numRows+col%numCols] = answer[i : i+1]
	}
	return grid
}

func IncrementClue(state *State, game *Game) error {
	number, direction := sampleBoard.Next(game.CurrentClue.Number, game.CurrentClue.Direction)
	return state.UpdateGameClue(game.ID, number, direction)
}
