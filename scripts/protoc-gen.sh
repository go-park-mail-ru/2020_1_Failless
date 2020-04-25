#!/usr/bin/env bash
protoc -I api/proto/ api/proto/auth.proto --go_out=plugins=grpc:api
