package main

import (
	"./settings"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"net/http"
	"strconv"
	"time"
)

func main() {
	serverSettings := settings.GetSettings()
	router := httptreemux.New()
	server := http.Server{
		Addr:              ":" + strconv.Itoa(serverSettings.Port),
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       time.Second * 10,
		WriteTimeout:      time.Second * 30,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
