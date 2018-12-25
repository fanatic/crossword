package main

import (
	"fmt"
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
	Guesses    map[string]Guess
}

func (state *State) CreateGame() (string, error) {
	id := GameID()
	w := GameState{
		ID:                   id,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CurrentClueNumber:    1,
		CurrentClueDirection: "across",
		CurrentClueExpiresAt: time.Now().Add(10 * time.Minute),
	}
	return id, state.db.Table("GameStates").Put(w).If("attribute_not_exists(ID)").Run()
}

func (state *State) CreatePlayer(gameID, name string) error {
	w := PlayerState{
		GameID:     gameID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		PlayerName: name,
	}
	return state.db.Table("PlayerStates").Put(w).If("attribute_not_exists(PlayerName)").Run()
}

func (state *State) CreateGuess(gameID, name, clueID, guess string) error {
	return state.db.Table("PlayerStates").
		Update("GameID", gameID).
		Range("PlayerName", name).
		Add("Guesses", map[string]Guess{clueID: Guess{Guess: guess, SubmittedAt: time.Now()}}).
		Set("UpdatedAt", time.Now()).
		Run()
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

	// table := db.Table("Games")

	// // put item
	// w := Game{UserID: 613, Time: time.Now(), Msg: "hello"}
	// err = table.Put(w).If("$ = ?", "Message", w.Msg).Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // get the same item
	// var result widget
	// err = table.Get("UserID", w.UserID).
	// 	Range("Time", dynamo.Equal, w.Time).
	// 	Filter("'Count' = ? AND $ = ?", w.Count, "Message", w.Msg). // placeholders in expressions
	// 	One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // get all items
	// var results []widget
	// err = table.Scan().All(&results)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return &State{
		db: db,
	}
}
