#!/bin/bash

go get -v ./...
GOOS=linux go build -o lambda_handler cmd/lambda_handler/main.go
zip handler.zip ./lambda_handler
aws lambda update-function-code \
  --function-name crossword \
  --zip-file fileb://./handler.zip 
