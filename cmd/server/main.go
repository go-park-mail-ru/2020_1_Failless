package main

import (
	"failless/configs/server"
	"failless/internal/pkg/settings"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if !settings.CheckSecretes(server.Secrets) {
		log.Fatal("Env variables don't set")
	}

	serverSettings := server.GetConfig()

	serve := http.Server{
		Addr:         serverSettings.Ip + ":" + strconv.Itoa(serverSettings.Port),
		Handler:      serverSettings.GetRouter(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	log.Println("server is running on " + strconv.Itoa(serverSettings.Port))
	//err := server.ListenAndServeTLS("./configs/ssl-bundle/bundle.crt", "./configs/ssl-bundle/private.key.pem")
	err := serve.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
