package server

import (
	"failless/configs/server"
	"failless/internal/pkg/logger"
	"failless/internal/pkg/settings"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Start() {
	file := logger.OpenLogFile("server")
	//defer file.Close()
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if ok := settings.CheckSecretes(server.Secrets); !ok {
		log.Println("Can't find variables ", server.Secrets)
		log.Fatal("Environment variables don't set")
	}
	serverSettings := server.GetConfig()
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
		log.Println(err)
	}
}
