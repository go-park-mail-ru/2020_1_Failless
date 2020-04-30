package router

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func OptionsReq(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("option request")
}

var promHandler = promhttp.Handler()


func MetricsHandler(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	promHandler.ServeHTTP(w, r)
}