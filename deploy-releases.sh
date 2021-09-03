#!/usr/bin/env sh

deploy() {
  aws s3 cp "builds/polarctl-${1}-${2}${3}" s3://polarctl/ --content-disposition "attachment; filename=\"polarctl${3}\""
}

deploy linux amd64
deploy linux arm64
deploy darwin amd64
deploy darwin arm64
deploy windows amd64 .exe
