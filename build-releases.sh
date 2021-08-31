#!/usr/bin/env sh

mkdir -p builds

export CGO_ENABLED=0

GOOS=linux   GOARCH=amd64 go build -ldflags "-X main.Version=${version}" -o builds/polarctl-linux-amd64

GOOS=linux   GOARCH=arm64 go build -ldflags "-X main.Version=${version}" -o builds/polarctl-linux-arm64

GOOS=darwin  GOARCH=amd64 go build -ldflags "-X main.Version=${version}" -o builds/polarctl-darwin-amd64

GOOS=darwin  GOARCH=arm64 go build -ldflags "-X main.Version=${version}" -o builds/polarctl-darwin-arm64

GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${version}" -o builds/polarctl-windows-amd64.exe
