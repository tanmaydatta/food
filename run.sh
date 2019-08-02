#!/usr/bin/env bash
git pull origin master
go get .
go build cmd/server/main.go
./main