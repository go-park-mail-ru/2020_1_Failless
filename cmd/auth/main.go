package auth

import (
	"failless/internal/app/auth"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	auth.Start()
}
