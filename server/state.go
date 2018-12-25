package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type State struct {
	db *dynamo.DB
}

func NewState() *State {
	creds := credentials.NewStaticCredentials("123", "123", "")
	db := dynamo.New(session.New(), &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-2"),
		Endpoint:    aws.String("http://localhost:8000"),
	})
	err := db.CreateTable("Games", Game{}).Run()
	if err != nil {
		log.Fatal(err)
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
