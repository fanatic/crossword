package main

import (
	"fmt"

	"github.com/akrylysov/algnhsa"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/fanatic/crossword/server"
)

func main() {
	err := xray.Configure(xray.Config{
		LogLevel:       "info", // default
		ServiceVersion: "1.2.3",
	})
	if err != nil {
		fmt.Println(err)
	}
	algnhsa.ListenAndServe(server.NewRouter(), nil)
}
