#!/usr/bin/env bash
git pull origin master
go build cmd/server/main.go
./main