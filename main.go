package main

import (
	"failless/server/routes"
	"failless/server/settings"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	serverSettings := settings.GetSettings()
	router := httptreemux.New()
	routes.AuthHandler(router)
	routes.SignUPHandler(router)
	routes.ProfileHandler(router)
	routes.TagHandler(router)
	routes.EventHandler(router)

	server := http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(serverSettings.Port),
		Handler:      router,
		TLSConfig:    nil,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	log.Println("server is running on " + strconv.Itoa(serverSettings.Port))
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
