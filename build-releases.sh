#!/usr/bin/env sh

mkdir -p builds
export CGO_ENABLED=0

build() {
  GOOS=${1} GOARCH=${2} go build \
    -ldflags "-X git.smith.care/smith/uc-phep/polar/polarctl/cmd.Version=${version}" \
    -o "builds/polarctl-${1}-${2}${3}"
}

build linux amd64
build linux arm64
build darwin amd64
build darwin arm64
build windows amd64 .exe
