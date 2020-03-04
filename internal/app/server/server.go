package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	sroutes "failless/internal/app/server/delivery"
)

func Start() {
	sroutes.ProfileHandler(router)
	sroutes.TagHandler(router)
	sroutes.EventHandler(router)

	server := http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(serverSettings.Port),
		Handler:      router,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	log.Println("server is running on " + strconv.Itoa(serverSettings.Port))
	//err := server.ListenAndServeTLS("./configs/ssl-bundle/bundle.crt", "./configs/ssl-bundle/private.key.pem")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
