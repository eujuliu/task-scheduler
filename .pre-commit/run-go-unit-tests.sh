#!/usr/bin/env bash

fail() {
  echo "unit tests failed"
  exit 1
}

FILES=$(go list ./... | grep -v /vendor/) || fail
go test -tags=unit -count=1 -short -v ${FILES} || fail
