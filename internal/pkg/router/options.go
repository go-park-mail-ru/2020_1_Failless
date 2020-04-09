package router

import (
	"log"
	"net/http"
)

func OptionsReq(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("option request")
}
