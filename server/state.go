package server

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type State struct {
	db *dynamo.DB
}

type GameState struct {
	ID                   string `dynamo:"ID,hash"`
	BoardID              string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CurrentClueNumber    int
	CurrentClueDirection string
	CurrentClueExpiresAt time.Time
}

type PlayerState struct {
	GameID     string `dynamo:"GameID,hash"`
	PlayerName string `dynamo:",range"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type GuessState struct {
	GameID       string `dynamo:"GameID,hash"`
	PlayerClueID string `dynamo:",range"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Guess        string
}

func PlayerClueID(name string, number int, direction string) string {
	return fmt.Sprintf("%s-%d-%s", name, number, direction)
}

func (g GuessState) PlayerName() string {
	return strings.Split(g.PlayerClueID, "-")[0]
}

func (g GuessState) ClueNumber() int {
	n, _ := strconv.Atoi(strings.Split(g.PlayerClueID, "-")[1])
	return n
}
func (g GuessState) ClueDirection() string {
	return strings.Split(g.PlayerClueID, "-")[2]
}

func (state *State) CreateGame(boardID string) (string, error) {
	id := GameID()
	w := GameState{
		ID:                   id,
		BoardID:              boardID,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CurrentClueNumber:    0,
		CurrentClueDirection: "across",
		CurrentClueExpiresAt: time.Now().Add(10 * time.Minute),
	}
	return id, state.db.Table("GameStates").Put(w).If("attribute_not_exists(ID)").Run()
}

func (state *State) UpdateGameClue(id string, number int, direction string) error {
	return state.db.Table("GameStates").
		Update("ID", id).
		Set("UpdatedAt", time.Now()).
		Set("CurrentClueNumber", number).
		Set("CurrentClueDirection", direction).
		Set("CurrentClueExpiresAt", time.Now().Add(30*time.Second)).
		Run()
}

func (state *State) CreatePlayer(gameID, name string) error {
	w := PlayerState{
		GameID:     gameID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		PlayerName: name,
	}
	err := state.db.Table("PlayerStates").Put(w).If("attribute_not_exists(PlayerName)").Run()
	if err == nil {
		return nil
	}

	return state.db.Table("PlayerStates").Update("GameID", gameID).
		Range("PlayerName", name).
		Set("UpdatedAt", time.Now()).
		Run()
}

func (state *State) CreateGuess(gameID, name string, clueNumber int, clueDirection, guess string) error {
	w := GuessState{
		GameID:       gameID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PlayerClueID: PlayerClueID(name, clueNumber, clueDirection),
		Guess:        guess,
	}
	err := state.db.Table("GuessStates").Put(w).If("attribute_not_exists(Guess)").Run()
	if err == nil {
		return nil
	}

	return state.db.Table("GuessStates").Update("GameID", gameID).
		Range("PlayerClueID", PlayerClueID(name, clueNumber, clueDirection)).
		Set("Guess", guess).
		Set("UpdatedAt", time.Now()).
		Run()
}

func (state *State) CreateBoardLayout(w BoardLayout) error {
	return state.db.Table("BoardLayouts").Put(w).Run()
}

func (state *State) GetBoardLayout(id string) (*BoardLayout, error) {
	var result BoardLayout
	err := state.db.Table("BoardLayouts").
		Get("ID", id).
		One(&result)
	return &result, err
}
func (state *State) GetBoardLayouts() ([]BoardLayout, error) {
	result := []BoardLayout{}
	err := state.db.Table("BoardLayouts").Scan().All(&result)
	return result, err
}

func (state *State) GetGame(gameID string) (*GameState, error) {
	var result GameState
	err := state.db.Table("GameStates").
		Get("ID", gameID).
		One(&result)
	return &result, err
}

func (state *State) GetPlayers(gameID string) ([]PlayerState, error) {
	var result []PlayerState
	err := state.db.Table("PlayerStates").
		Get("GameID", gameID).
		All(&result)
	return result, err
}

func (state *State) GetGuesses(gameID string) ([]GuessState, error) {
	var result []GuessState
	err := state.db.Table("GuessStates").
		Get("GameID", gameID).
		All(&result)
	return result, err
}

func NewState() *State {
	creds := credentials.NewStaticCredentials("123", "123", "")
	db := dynamo.New(session.New(), &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-2"),
		Endpoint:    aws.String("http://localhost:8000"),
	})
	err := db.CreateTable("GameStates", GameState{}).Run()
	if err != nil {
		fmt.Println(err)
	}

	err = db.CreateTable("PlayerStates", PlayerState{}).Run()
	if err != nil {
		fmt.Println(err)
	}

	err = db.CreateTable("GuessStates", GuessState{}).Run()
	if err != nil {
		fmt.Println(err)
	}

	err = db.CreateTable("BoardLayouts", BoardLayout{}).Run()
	if err != nil {
		fmt.Println(err)
	}

	return &State{
		db: db,
	}
}
