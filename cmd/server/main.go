package main

import (
	"failless/internal/app/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	server.Start()
}
