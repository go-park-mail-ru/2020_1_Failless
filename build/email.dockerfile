FROM golang:1.13-alpine3.11

MAINTAINER Failless

WORKDIR /home/eventum
COPY . .
RUN apk add --no-cache git wget unzip
RUN apk add --no-cache git wget unzip protobuf
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promauto
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get -u github.com/mailru/easyjson/...

RUN easyjson -all internal/pkg/models/models.go
RUN easyjson -all internal/pkg/forms/form_*.go
RUN easyjson internal/pkg/security/security.go

RUN go build -o bin/email ./cmd/email/main.go
