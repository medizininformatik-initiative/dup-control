#!/usr/bin/env sh

aws s3 cp "builds/VERSION" s3://dupctl/

deploy() {
  aws s3 cp "builds/dupctl-${1}-${2}${3}" s3://dupctl/ --content-disposition "attachment; filename=\"dupctl${3}\""
}

deploy linux amd64
deploy linux arm64
deploy darwin amd64
deploy darwin arm64
deploy windows amd64 .exe
