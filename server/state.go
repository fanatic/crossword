package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
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
	CurrentRound         int
	LastClueNumber       int
	LastClueDirection    string
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

func PlayerClueIDRound(name string, number int, direction string, round int) string {
	return fmt.Sprintf("%s-%d-%s-%d", name, number, direction, round)
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
func (g GuessState) Round() int {
	n, _ := strconv.Atoi(strings.Split(g.PlayerClueID, "-")[3])
	return n
}

func (state *State) CreateGame(ctx context.Context, boardID string) (string, error) {
	id := GameID()
	w := GameState{
		ID:                   id,
		BoardID:              boardID,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CurrentClueNumber:    0,
		CurrentClueDirection: "across",
		CurrentClueExpiresAt: time.Now().Add(10 * time.Minute),
		LastClueNumber:       0,
		LastClueDirection:    "across",
		CurrentRound:         0,
	}
	return id, state.db.Table("GameStates").Put(w).If("attribute_not_exists(ID)").RunWithContext(ctx)
}

func (state *State) UpdateGameClue(ctx context.Context, id string, number int, direction string, lastNumber int, lastDirection string, round int) error {
	return state.db.Table("GameStates").
		Update("ID", id).
		Set("UpdatedAt", time.Now()).
		Set("CurrentClueNumber", number).
		Set("CurrentClueDirection", direction).
		Set("CurrentClueExpiresAt", time.Now().Add(30*time.Second)).
		Set("LastClueNumber", lastNumber).
		Set("LastClueDirection", direction).
		Set("CurrentRound", round).
		RunWithContext(ctx)
}

func (state *State) CreatePlayer(ctx context.Context, gameID, name string) error {
	w := PlayerState{
		GameID:     gameID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		PlayerName: name,
	}
	err := state.db.Table("PlayerStates").Put(w).If("attribute_not_exists(PlayerName)").RunWithContext(ctx)
	if err == nil {
		return nil
	}

	return state.db.Table("PlayerStates").Update("GameID", gameID).
		Range("PlayerName", name).
		Set("UpdatedAt", time.Now()).
		RunWithContext(ctx)
}

func (state *State) CreateGuess(ctx context.Context, gameID, name string, clueNumber int, clueDirection, guess string, round int) error {
	w := GuessState{
		GameID:       gameID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PlayerClueID: PlayerClueIDRound(name, clueNumber, clueDirection, round),
		Guess:        guess,
	}
	err := state.db.Table("GuessStates").Put(w).If("attribute_not_exists(Guess)").RunWithContext(ctx)
	if err == nil {
		return nil
	}

	return state.db.Table("GuessStates").Update("GameID", gameID).
		Range("PlayerClueID", PlayerClueIDRound(name, clueNumber, clueDirection, round)).
		Set("Guess", guess).
		Set("UpdatedAt", time.Now()).
		RunWithContext(ctx)
}

func (state *State) CreateBoardLayout(w BoardLayout) error {
	return state.db.Table("BoardLayouts").Put(w).Run()
}

func (state *State) GetBoardLayout(ctx context.Context, id string) (*BoardLayout, error) {
	var result BoardLayout
	err := state.db.Table("BoardLayouts").
		Get("ID", id).
		OneWithContext(ctx, &result)
	return &result, err
}
func (state *State) GetBoardLayouts(ctx context.Context) ([]BoardLayout, error) {
	result := []BoardLayout{}
	err := state.db.Table("BoardLayouts").Scan().AllWithContext(ctx, &result)
	return result, err
}

func (state *State) GetGame(ctx context.Context, gameID string) (*GameState, error) {
	var result GameState
	err := state.db.Table("GameStates").
		Get("ID", gameID).
		OneWithContext(ctx, &result)
	return &result, err
}

func (state *State) GetPlayers(ctx context.Context, gameID string) ([]PlayerState, error) {
	var result []PlayerState
	err := state.db.Table("PlayerStates").
		Get("GameID", gameID).
		AllWithContext(ctx, &result)
	return result, err
}

func (state *State) GetGuesses(ctx context.Context, gameID string) ([]GuessState, error) {
	var result []GuessState
	err := state.db.Table("GuessStates").
		Get("GameID", gameID).
		AllWithContext(ctx, &result)
	return result, err
}

func NewState() *State {
	//creds := credentials.NewStaticCredentials("123", "123", "")
	db := dynamo.New(session.New(), &aws.Config{
		//Credentials: creds,
		Region: aws.String("us-east-2"),
		//Endpoint:    aws.String("http://localhost:8000"),
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
