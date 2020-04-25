FROM golang:1.13-stretch AS lang
MAINTAINER Failless

WORKDIR /home/eventum
COPY . .
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN go build -o bin/eventum ./cmd/server/main.go

EXPOSE 3001

CMD ./bin/eventum
