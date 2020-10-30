#!/bin/sh

env GOOS=linux GOARCH=amd64 CGO_ENABLED=1

cd /go/src/github.com/xiaof-github/kafdrop-go
rm -f kafdrop
go build -ldflags "-linkmode external -extldflags -static" -o kafdrop main.go

echo "install finish."