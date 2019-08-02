#!/usr/bin/env bash
git pull origin master
cd cmd/server
go get .
env GOOS=linux GOARCH=amd64 go build -v .
#env GOOS=darwin go build -v .
./server