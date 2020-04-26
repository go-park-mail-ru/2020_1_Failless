#!/usr/bin/env bash

protoc -I api/proto/ api/proto/auth.proto --go_out=plugins=grpc:api
go build -o bin/auth ./cmd/auth
go build -o bin/server ./cmd/server
go build -o bin/chat ./cmd/chat

