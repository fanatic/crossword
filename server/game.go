package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
)

type Game struct {
	ID             string       `json:"id"`
	BoardLayout    *BoardLayout `json:"layout"`
	LastClue       Clue         `json:"last_clue"`
	CurrentClue    Clue         `json:"current_clue"`
	CurrentPlayers []Player     `json:"current_players"`
	CurrentRound   int          `json:"-"`
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
	IsCorrect   bool      `json:"-"`
	LatestRound int       `json:"-"`
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

func WaitingOnPlayers(currentPlayers []Player, guesses []Guess, currentRound int) []string {
	players := map[string]bool{}
	for _, p := range currentPlayers {
		if p.Active {
			players[p.Name] = true
		}
	}
	for _, g := range guesses {
		if currentRound == g.LatestRound {
			delete(players, g.Player.Name)
		}
	}

	waitingOnPlayers := []string{}
	for player := range players {
		waitingOnPlayers = append(waitingOnPlayers, player)
	}
	return waitingOnPlayers
}

func FetchGame(ctx context.Context, state *State, gameID string) (*Game, error) {
	var result *Game
	var err error
	deadline, _ := ctx.Deadline()
	deadline = deadline.Add(-100 * time.Millisecond)
	timeoutChannel := time.After(time.Until(deadline))

	done := make(chan struct{})
	go func() {
		result, err = fetchGameReal(ctx, state, gameID)
		close(done)
	}()

	select {
	case <-timeoutChannel:
		return nil, fmt.Errorf("Timed out")

	case <-done:
		return result, err
	}
}

func fetchGameReal(ctx context.Context, state *State, gameID string) (*Game, error) {
	game, err := state.GetGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	board, err := state.GetBoardLayout(ctx, game.BoardID)
	if err != nil {
		return nil, err
	}

	players, err := state.GetPlayers(ctx, gameID)
	if err != nil {
		return nil, err
	}

	guesses, err := state.GetGuesses(ctx, gameID)
	if err != nil {
		return nil, err
	}

	clueGuesses := map[string][]Guess{}
	playerGuesses := map[string][]Guess{}
	xray.Capture(ctx, "process-guesses", func(ctx1 context.Context) error {
		for _, guess := range guesses {
			g := Guess{
				Clue:        Clue{Number: guess.ClueNumber(), Direction: guess.ClueDirection()},
				Player:      Player{Name: guess.PlayerName()},
				Guess:       guess.Guess,
				SubmittedAt: guess.UpdatedAt,
				LatestRound: guess.Round(),
			}
			g.Score, g.IsCorrect = board.CalculateScore(g.Guess, g.Clue.Number, g.Clue.Direction)

			clueID := fmt.Sprintf("%d-%s", g.Clue.Number, g.Clue.Direction)
			clueGuesses[clueID] = AddOrSetGuess(board, g, clueGuesses[clueID])
			playerGuesses[guess.PlayerName()] = AddOrSetGuess(board, g, playerGuesses[guess.PlayerName()])
		}
		return nil
	})

	grid := InitializeGrid(board.Size.Rows, board.Size.Cols, board.Grid, board.Gridnums)
	xray.Capture(ctx, "process-grid", func(ctx1 context.Context) error {
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
		return nil
	})

	g := Game{
		ID:             game.ID,
		BoardLayout:    board,
		CurrentClue:    board.GetClue(game.CurrentClueNumber, game.CurrentClueDirection, &game.CurrentClueExpiresAt),
		LastClue:       board.GetClue(game.LastClueNumber, game.LastClueDirection, nil),
		CurrentPlayers: []Player{},
	}

	lastClueID := fmt.Sprintf("%d-%s", g.LastClue.Number, g.LastClue.Direction)
	g.LastClue.Guesses = clueGuesses[lastClueID]
	g.LastClue.Answer = ""

	xray.Capture(ctx, "process-players", func(ctx1 context.Context) error {
		for _, player := range players {
			p := Player{
				Name:         player.PlayerName,
				CurrentScore: 0,
				Guesses:      playerGuesses[player.PlayerName],
				Active:       true,
			}
			for _, g := range p.Guesses {
				p.CurrentScore += g.Score
			}
			g.CurrentPlayers = append(g.CurrentPlayers, p)
		}
		return nil
	})

	currentClueID := fmt.Sprintf("%d-%s", g.CurrentClue.Number, g.CurrentClue.Direction)
	g.CurrentClue.WaitingOnPlayers = WaitingOnPlayers(g.CurrentPlayers, clueGuesses[currentClueID], g.CurrentRound)
	g.CurrentClue.Answer = board.MaskAnswer(g.CurrentClue.Answer, g.CurrentClue.Number, g.CurrentClue.Direction, grid)

	g.BoardLayout.Grid = grid
	g.BoardLayout.Answers = nil
	return &g, nil
}

func AddOrSetGuess(board *BoardLayout, g Guess, guesses []Guess) []Guess {
	for _, guess := range guesses {
		if guess.Player.Name == g.Player.Name && guess.Clue.Number == g.Clue.Number && guess.Clue.Direction == g.Clue.Direction {
			// Merge new round with previous round
			if !guess.IsCorrect && g.IsCorrect {
				guess.IsCorrect = true
				guess.Score = g.Score
			} else if !guess.IsCorrect {
				if guess.LatestRound < g.LatestRound && !guess.IsCorrect {
					guess.Guess = g.Guess
					guess.LatestRound = g.LatestRound
				}
			}
			return guesses
		}
	}
	return append(guesses, g)
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

func IncrementClue(ctx context.Context, state *State, game *Game) error {
	board, err := state.GetBoardLayout(ctx, game.BoardLayout.ID)
	if err != nil {
		return err
	}

	round := game.CurrentRound
	nextRound := false
	number, direction := game.BoardLayout.Next(game.CurrentClue.Number, game.CurrentClue.Direction)
	for i := 0; i <= len(board.Answers.Across)-1+len(board.Answers.Down)-1; i++ {
		// Increment round
		if game.CurrentClue.Direction == "down" && direction == "across" {
			nextRound = true
		}
		// Skip over answered clues
		if !ClueAnswered(state, game, board, number, direction) {
			break
		}
		number, direction = game.BoardLayout.Next(number, direction)
	}
	if nextRound {
		round++
	}

	return state.UpdateGameClue(ctx, game.ID, number, direction, game.CurrentClue.Number, game.CurrentClue.Direction, round)
}

func ClueAnswered(state *State, game *Game, board *BoardLayout, clueNumber int, clueDirection string) bool {
	correctAnswer := board.CorrectAnswer(clueNumber, clueDirection)
	gridAnswer := board.MaskAnswer(correctAnswer, clueNumber, clueDirection, game.BoardLayout.Grid)

	return correctAnswer == gridAnswer
}

func CheckTimers(ctx context.Context, state *State, game *Game) error {
	if time.Since(*game.CurrentClue.ExpiresAt) < 0 {
		return nil
	}

	return IncrementClue(ctx, state, game)
}
