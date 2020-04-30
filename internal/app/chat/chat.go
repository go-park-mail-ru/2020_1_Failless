package chat

import (
	"failless/configs/chat"
	"failless/internal/pkg/logger"
	"failless/internal/pkg/settings"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Start() {
	file := logger.OpenLogFile("chat")
	defer file.Close()

	if ok := settings.CheckSecretes(chat.Secrets); !ok {
		log.Println("Can't find variables ", chat.Secrets)
		log.Fatal("Environment variables don't set")
	}
	serverSettings := chat.GetConfig()
	serve := http.Server{
		Addr:              serverSettings.Ip + ":" + strconv.Itoa(serverSettings.Port),
		Handler:           serverSettings.GetRouter(),
		ReadTimeout:       time.Second * 10,
		WriteTimeout:      time.Second * 30,
		ReadHeaderTimeout: time.Second * 30,
		IdleTimeout:       time.Second * 120,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
	}

	log.Println("chat server is running on " + strconv.Itoa(serverSettings.Port))
	err := serve.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
