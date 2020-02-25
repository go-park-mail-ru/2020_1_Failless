package main

import (
	"./routes"
	"./settings"
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
	routes.SignInHandler(router)
	routes.SignUPHandler(router)
	server := http.Server{
		Addr:         ":" + strconv.Itoa(serverSettings.Port),
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
