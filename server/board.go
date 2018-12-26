package server

import (
	"html"
	"strconv"
	"strings"
	"time"
)

type BoardLayout struct {
	Admin   bool `json:"admin"`
	Answers *struct {
		Across []string `json:"across"`
		Down   []string `json:"down"`
	} `json:"answers"`
	Author   string      `json:"author"`
	Autowrap interface{} `json:"autowrap"`
	Bbars    interface{} `json:"bbars"`
	Circles  interface{} `json:"circles"`
	Clues    struct {
		Across []string `json:"across"`
		Down   []string `json:"down"`
	} `json:"clues"`
	Code            interface{} `json:"code"`
	Copyright       string      `json:"copyright"`
	Date            string      `json:"date"`
	Dow             string      `json:"dow"`
	Downmap         interface{} `json:"downmap"`
	Editor          string      `json:"editor"`
	Grid            []string    `json:"grid"`
	Gridnums        []int       `json:"gridnums"`
	Hastitle        bool        `json:"hastitle"`
	Hold            interface{} `json:"hold"`
	ID              string      `json:"id" dynamo:"ID,hash"`
	ID2             interface{} `json:"id2"`
	Interpretcolors interface{} `json:"interpretcolors"`
	Jnotes          interface{} `json:"jnotes"`
	Key             interface{} `json:"key"`
	Mini            interface{} `json:"mini"`
	Notepad         interface{} `json:"notepad"`
	Publisher       string      `json:"publisher"`
	Rbars           interface{} `json:"rbars"`
	Shadecircles    interface{} `json:"shadecircles"`
	Size            struct {
		Cols int `json:"cols"`
		Rows int `json:"rows"`
	} `json:"size"`
	Title   string      `json:"title"`
	Track   interface{} `json:"track"`
	Type    interface{} `json:"type"`
	Uniclue bool        `json:"uniclue"`
	Valid   bool        `json:"valid"`
}

func (b BoardLayout) ClueLabel(number int, direction string) int {
	clues := b.Clues.Down
	if direction == "across" {
		clues = b.Clues.Across
	}
	label := strings.Split(clues[number], ".")[0]
	labelInt, _ := strconv.Atoi(label)
	return labelInt
}

func (b BoardLayout) ClueLabelPosition(clueLabel int) (int, int) {
	for idx, gridnum := range b.Gridnums {
		if gridnum == clueLabel {
			row := idx / b.Size.Rows
			col := idx % b.Size.Cols
			return row, col
		}
	}
	return -1, -1
}

func (b BoardLayout) GetLastClue(number int, direction string) Clue {
	prevNum, prevDirection := b.Prev(number, direction)
	return b.GetClue(prevNum, prevDirection, nil)
}

func (b BoardLayout) GetClue(number int, direction string, expiresAt *time.Time) Clue {
	clues := b.Clues.Down
	answers := b.Answers.Down
	if direction == "across" {
		clues = b.Clues.Across
		answers = b.Answers.Across
	}
	return Clue{
		Number:      number,
		Direction:   direction,
		Description: html.UnescapeString(clues[number]),
		Answer:      answers[number],
		ExpiresAt:   expiresAt,
	}
}

func (b BoardLayout) Next(number int, direction string) (int, string) {
	if direction == "across" {
		if len(b.Clues.Across) == number+1 {
			return 0, "down"
		}
	} else {
		if len(b.Clues.Down) == number+1 {
			return 0, "across"
		}
	}
	return number + 1, direction
}

func (b BoardLayout) Prev(number int, direction string) (int, string) {
	if number == 0 {
		if direction == "across" {
			return len(b.Clues.Down) - 1, "down"
		} else {
			return len(b.Clues.Across) - 1, "across"
		}
	}
	return number - 1, direction
}

func (b *BoardLayout) CorrectAnswer(clueNumber int, clueDirection string) string {
	if clueDirection == "across" {
		return b.Answers.Across[clueNumber]
	}
	return b.Answers.Down[clueNumber]
}

func (b *BoardLayout) CalculateScore(guess Guess) int {
	if strings.ToUpper(guess.Guess) != b.CorrectAnswer(guess.Clue.Number, guess.Clue.Direction) {
		return 0
	}

	return 100
}

func (b *BoardLayout) MaskAnswer(answer string, currentClueNumber int, currentClueDirection string, grid []string) string {
	clueLabel := b.ClueLabel(currentClueNumber, currentClueDirection)
	row, col := b.ClueLabelPosition(clueLabel)
	numRows, numCols := b.Size.Rows, b.Size.Cols

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
