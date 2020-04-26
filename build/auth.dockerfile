FROM golang:1.13-alpine3.11

MAINTAINER Failless

WORKDIR /home/eventum
COPY . .
RUN apk add --no-cache git wget unzip
ENV PROTOCV 3.11.4
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOCV/protoc-$PROTOCV-linux-x86_64.zip -P /tmp/ &&\
RUN unzip /tmp/protoc-$PROTOCV-linux-x86_64.zip -d /tmp/
RUN cp /tmp/protoc-$PROTOCV-linux-x86_64/bin/* /usr/bin/.

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN /usr/bin/protoc -I /home/eventum/api/proto/ /home/eventum/api/proto/auth.proto --go_out=plugins=grpc:/home/eventum/api
RUN go build -o bin/auth ./cmd/auth/main.go
