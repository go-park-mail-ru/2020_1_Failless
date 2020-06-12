package main

import (
	"failless/internal/app/email"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	email.Start()
}
