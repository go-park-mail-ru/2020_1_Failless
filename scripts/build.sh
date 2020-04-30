#!/usr/bin/env bash

protoc -I api/proto/ api/proto/auth.proto --go_out=plugins=grpc:api
echo 'proto files were generated'
rm internal/pkg/forms/*_easyjson.go
easyjson -all internal/pkg/models/models.go
echo 'easyjson models were generated'
easyjson -all internal/pkg/forms/form_*.go
echo 'easyjson forms were generated'
easyjson internal/pkg/security/security.go
echo 'easyjson security file was generated'

go build -o bin/auth ./cmd/auth
go build -o bin/server ./cmd/server
go build -o bin/chat ./cmd/chat

