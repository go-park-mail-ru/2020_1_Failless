#!/usr/bin/env bash

go generate ./...
go test `go list all | grep "failless" | grep -v "mocks"` -coverprofile=coverage.out.tmp -cover ./...
cat coverage.out.tmp | grep -v ".*_easyjson.go.*" > coverage.out
go tool cover -func=coverage.out
rm coverage.out.tmp