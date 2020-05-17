FROM golang:1.13-alpine3.11

MAINTAINER Failless

WORKDIR /home/eventum
COPY . .
RUN apk add --no-cache git wget unzip protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN go get google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN protoc -I /home/eventum/api/proto/ /home/eventum/api/proto/auth.proto --go_out=plugins=grpc:/home/eventum/api

RUN go get -u github.com/mailru/easyjson/...
RUN easyjson -all internal/pkg/models/models.go
RUN easyjson -all internal/pkg/forms/form_*.go
RUN easyjson internal/pkg/security/security.go

RUN go build -o bin/auth ./cmd/auth/main.go
