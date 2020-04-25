package auth

import (
	"failless/configs/auth"
	"failless/internal/pkg/logger"
	"failless/internal/pkg/settings"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Start() {
	file := logger.OpenLogFile("auth")
	defer file.Close()

	if ok := settings.CheckSecretes(auth.Secrets); !ok {
		log.Println("Can't find variables ", auth.Secrets)
		log.Fatal("Environment variables don't set")
	}
	serverSettings := auth.GetConfig()
	serve := http.Server{
		Addr:              serverSettings.Ip + ":" + strconv.Itoa(serverSettings.Port),
		Handler:           serverSettings.GetRouter(),
		ReadTimeout:       time.Second * 10,
		WriteTimeout:      time.Second * 30,
		ReadHeaderTimeout: time.Second * 30,
		IdleTimeout:       time.Second * 120,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes, // TODO: ask what is the best value
	}

	log.Println("server is running on " + strconv.Itoa(serverSettings.Port))
	err := serve.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
