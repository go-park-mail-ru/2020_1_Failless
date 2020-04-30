FROM golang:1.13-alpine3.11

MAINTAINER Failless

WORKDIR /home/eventum
COPY . .
RUN apk add --no-cache git wget unzip
RUN apk add --no-cache git wget unzip protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN go get google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promauto
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get -u github.com/mailru/easyjson/...

RUN protoc -I /home/eventum/api/proto/ /home/eventum/api/proto/auth.proto --go_out=plugins=grpc:/home/eventum/api
RUN go build -o bin/chat ./cmd/chat/main.go
