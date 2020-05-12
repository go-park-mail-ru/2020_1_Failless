TARGET_SERVER = server
TARGET_AUTH = auth
TARGET_CHAT = chat
GO = go

# PID file will keep the process id of the server
PID1 := /tmp/.$(TARGET_SERVER).pid
PID2 := /tmp/.$(TARGET_AUTH).pid
PID3 := /tmp/.$(TARGET_CHAT).pid

GOFILES := $(PWD)/..


# Make is verbose in Linux. Make it silent.
#MAKEFLAGS += --silent

.PHONY: help
all: help
help: Makefile
	@echo
	@echo "Choose a command run in Eventum-Backend:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo


all: go-get generate compile

## Start: Run development mode.
start: start-auth start-server start-chat
## stop: Stop development mode.
stop: stop-server stop-auth stop-chat

clean: clean-bin clean-log clean-json

generate: protoc easyjson

compile: build-server build-auth build-chat

start-server:
	@echo "  >  $(TARGET_SERVER) is available at localhost:$(SERVER_PORT)"
	@-$(GOFILES)/bin/server 2>&1 & echo $$! > $(PID1)
	@cat $(PID1) | sed "/^/s/^/  \>  PID: /"

start-auth:
	@echo "  >  $(TARGET_AUTH) is available at localhost:$(AUTH_PORT)"
	@-$(GOFILES)/bin/auth 2>&1 & echo $$! > $(PID2)
	@cat $(PID2) | sed "/^/s/^/  \>  PID: /"

start-chat:
	@echo "  >  $(TARGET_CHAT) is available at localhost:$(CHAT_PORT)"
	@-$(GOFILES)/bin/chat 2>&1 & echo $$! > $(PID3)
	@cat $(PID3) | sed "/^/s/^/  \>  PID: /"

stop-server:
	@-touch $(PID1)
	@-kill `cat $(PID1)` 2> /dev/null || true
	@-rm $(PID1)

stop-auth:
	@-touch $(PID2)
	@-kill `cat $(PID2)` 2> /dev/null || true
	@-rm $(PID2)

stop-chat:
	@-touch $(PID3)
	@-kill `cat $(PID3)` 2> /dev/null || true
	@-rm $(PID3)

build-server:
	@echo "  >  Building binary for $(TARGET_SERVER)..."
	$(GO) build $(FLAGS) -o $(GOFILES)/bin/server $(GOFILES)/cmd/server

build-auth:
	@echo "  >  Building binary for $(TARGET_AUTH)..."
	$(GO) build $(FLAGS) -o $(GOFILES)/bin/auth $(GOFILES)/cmd/auth

build-chat:
	@echo "  >  Building binary for $(TARGET_CHAT)..."
	$(GO) build $(FLAGS) -o $(GOFILES)/bin/chat $(GOFILES)/cmd/chat

easyjson:
	@echo "  >  Generating easyjson files..."
	easyjson -all $(GOFILES)/internal/pkg/models/models.go
	easyjson -all $(GOFILES)/internal/pkg/forms/form_*.go
	easyjson $(GOFILES)/internal/pkg/security/security.go

protoc:
	@echo "  >  Generating proto files..."
	protoc -I $(GOFILES)/api/proto/ $(GOFILES)/api/proto/auth.proto --go_out=plugins=grpc:../api

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	$(GO) get -t

clean-bin:
	@echo "  >  Cleaning build cache"
	rm -rf $(GOFILES)/bin

clean-log:
	@echo "  >  Cleaning app logs"
	rm $(GOFILES)/deployments/*.log
	rm $(GOFILES)/*.log

clean-json:
	rm $(GOFILES)/internal/pkg/forms/*_easyjson.go
	rm $(GOFILES)/internal/pkg/forms/easyjson-*.go

