package main

import (
	"failless/internal/app/chat"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	chat.Start()
}
