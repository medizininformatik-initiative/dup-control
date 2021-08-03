#!/usr/bin/env sh

mkdir -p builds

export CGO_ENABLED=0

GOOS=linux   GOARCH=amd64 go build
tar czf builds/polarctl-${version}-linux-amd64.tar.gz polarctl
rm polarctl

GOOS=linux   GOARCH=arm64 go build
tar czf builds/polarctl-${version}-linux-arm64.tar.gz polarctl
rm polarctl

GOOS=darwin  GOARCH=amd64 go build
tar czf builds/polarctl-${version}-darwin-amd64.tar.gz polarctl
rm polarctl

GOOS=darwin  GOARCH=arm64 go build
tar czf builds/polarctl-${version}-darwin-arm64.tar.gz polarctl
rm polarctl

GOOS=windows GOARCH=amd64 go build
zip -q builds/polarctl-${version}-windows-amd64.zip polarctl.exe
rm polarctl.exe
