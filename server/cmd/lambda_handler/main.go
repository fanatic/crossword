package main

import (
	"github.com/akrylysov/algnhsa"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/fanatic/crossword/server"
)

func main() {
	xray.Configure(xray.Config{
		LogLevel:       "info", // default
		ServiceVersion: "1.2.3",
	})
	algnhsa.ListenAndServe(server.NewRouter(), nil)
}
