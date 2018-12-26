package server

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	ID             string       `json:"id"`
	BoardLayout    *BoardLayout `json:"layout"`
	LastClue       Clue         `json:"last_clue"`
	CurrentClue    Clue         `json:"current_clue"`
	CurrentPlayers []Player     `json:"current_players"`
}

type Player struct {
	Name         string  `json:"name"`
	CurrentScore int     `json:"current_score"`
	Guesses      []Guess `json:"-"`
	Active       bool    `json:"active"`
}

type Guess struct {
	Player      Player    `json:"player" dynamo:"-"`
	Clue        Clue      `json:"-" dynamo:"-"`
	Guess       string    `json:"guess"`
	Score       int       `json:"score"`
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
		if p.Active {
			players[p.Name] = true
		}
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

	board, err := state.GetBoardLayout(game.BoardID)
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
		g.Score = board.CalculateScore(g)
		clueID := fmt.Sprintf("%d-%s", guess.ClueNumber(), guess.ClueDirection())
		clueGuesses[clueID] = append(clueGuesses[clueID], g)
		playerGuesses[guess.PlayerName()] = append(playerGuesses[guess.PlayerName()], g)
	}

	grid := InitializeGrid(board.Size.Rows, board.Size.Cols, board.Grid, board.Gridnums)
	for clueID, guesses := range clueGuesses {
		n, _ := strconv.Atoi(strings.Split(clueID, "-")[0])
		d := strings.Split(clueID, "-")[1]

		// Skip evaluating this clue if it's the current one
		if game.CurrentClueNumber == n && game.CurrentClueDirection == d {
			continue
		}

		for _, guess := range guesses {
			// If guess is correct, fill in grid with answer
			answer := strings.ToUpper(guess.Guess)
			correctAnswer := board.CorrectAnswer(n, d)
			if correctAnswer == answer {
				clueLabel := board.ClueLabel(n, d)
				row, col := board.ClueLabelPosition(clueLabel)
				if d == "across" {
					grid = FillInGridAnswerAcross(row, col, board.Size.Rows, board.Size.Cols, answer, grid)
				} else {
					grid = FillInGridAnswerDown(row, col, board.Size.Rows, board.Size.Cols, answer, grid)
				}
			}
		}
	}

	g := Game{
		ID:             game.ID,
		BoardLayout:    board,
		CurrentClue:    board.GetClue(game.CurrentClueNumber, game.CurrentClueDirection, &game.CurrentClueExpiresAt),
		LastClue:       board.GetLastClue(game.CurrentClueNumber, game.CurrentClueDirection),
		CurrentPlayers: []Player{},
	}

	lastClueID := fmt.Sprintf("%d-%s", g.LastClue.Number, g.LastClue.Direction)
	g.LastClue.Guesses = clueGuesses[lastClueID]
	g.LastClue.Answer = ""

	for _, player := range players {
		p := Player{
			Name:         player.PlayerName,
			CurrentScore: 0,
			Guesses:      playerGuesses[player.PlayerName],
			Active:       false,
		}
		for _, guess := range p.Guesses {
			p.CurrentScore += board.CalculateScore(guess)
		}
		for _, guess := range g.LastClue.Guesses {
			if guess.Player.Name == player.PlayerName {
				p.Active = true
			}
		}
		g.CurrentPlayers = append(g.CurrentPlayers, p)
	}

	currentClueID := fmt.Sprintf("%d-%s", g.CurrentClue.Number, g.CurrentClue.Direction)
	g.CurrentClue.WaitingOnPlayers = WaitingOnPlayers(g.CurrentPlayers, clueGuesses[currentClueID])
	g.CurrentClue.Answer = board.MaskAnswer(g.CurrentClue.Answer, g.CurrentClue.Number, g.CurrentClue.Direction, grid)

	g.BoardLayout.Grid = grid
	g.BoardLayout.Answers = nil
	return &g, nil
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
	board, err := state.GetBoardLayout(game.BoardLayout.ID)
	if err != nil {
		return err
	}

	number, direction := game.BoardLayout.Next(game.CurrentClue.Number, game.CurrentClue.Direction)
	for i := 0; i <= len(board.Answers.Across)-1+len(board.Answers.Down)-1; i++ {
		// Skip over answered clues
		if !ClueAnswered(state, game, board, number, direction) {
			break
		}
		number, direction = game.BoardLayout.Next(number, direction)
	}

	return state.UpdateGameClue(game.ID, number, direction)
}

func ClueAnswered(state *State, game *Game, board *BoardLayout, clueNumber int, clueDirection string) bool {
	correctAnswer := board.CorrectAnswer(clueNumber, clueDirection)
	gridAnswer := board.MaskAnswer(correctAnswer, clueNumber, clueDirection, game.BoardLayout.Grid)

	return correctAnswer == gridAnswer
}

func CheckTimers(state *State, game *Game) error {
	if time.Since(*game.CurrentClue.ExpiresAt) < 0 {
		return nil
	}

	return IncrementClue(state, game)
}
