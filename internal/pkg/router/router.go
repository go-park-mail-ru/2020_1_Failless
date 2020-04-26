package router

import (
	"failless/internal/pkg/middleware"
	"failless/internal/pkg/settings"
	"github.com/dimfeld/httptreemux"
	"log"
)

// Parse route map and return configured Router
func InitRouter(s *settings.ServerSettings, router *httptreemux.TreeMux) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Error was occurred", r)
		}
	}()

	var optionsHandler settings.HandlerFunc = nil
	for key, list := range s.Routes {
		for _, pack := range list {
			handler := pack.Handler
			if pack.Type != "OPTIONS" {
				handler = middleware.SetCSRF(handler)
				// TODO: check is csrf token work
				if pack.CSRF {
					handler = middleware.CheckCSRF(handler)
				}
			}

			if pack.CORS {
				s.Secure.CORSMap[pack.Type] = struct{}{}
				handler = middleware.CORS(handler)
			}

			if pack.AuthRequired {
				handler = middleware.Auth(handler)
			}

			switch pack.Type {
			case "GET":
				(*router).GET(key, httptreemux.HandlerFunc(handler))
			case "PUT":
				(*router).PUT(key, httptreemux.HandlerFunc(handler))
			case "POST":
				(*router).POST(key, httptreemux.HandlerFunc(handler))
			case "DELETE":
				(*router).DELETE(key, httptreemux.HandlerFunc(handler))
			case "OPTIONS":
				optionsHandler = handler
			}

		}
	}

	if optionsHandler != nil {
		for key, _ := range s.Routes {
			(*router).OPTIONS(key, httptreemux.HandlerFunc(optionsHandler))
		}
	}
	// generate "GET, POST, OPTIONS, HEAD, PUT" string
	for key, _ := range s.Secure.CORSMap {
		s.Secure.CORSMethods += key + ", "
	}

	// remove extra comma
	s.Secure.CORSMethods = s.Secure.CORSMethods[:len(s.Secure.CORSMethods)-2]
	s.Router = router
}



func InitWSRoter(s *settings.ServerSettings, router *httptreemux.TreeMux) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Error was occurred", r)
		}
	}()

	var optionsHandler settings.HandlerFunc = nil
	for key, list := range s.Routes {
		for _, pack := range list {
			handler := pack.Handler
			if pack.Type != "OPTIONS" {
				handler = middleware.SetCSRF(handler)
				// TODO: check is csrf token work
				if pack.CSRF {
					handler = middleware.CheckCSRF(handler)
				}
			}

			if pack.CORS {
				s.Secure.CORSMap[pack.Type] = struct{}{}
				handler = middleware.CORS(handler)
			}

			if pack.AuthRequired {
				handler = middleware.Auth(handler)
			}

			switch pack.Type {
			case "GET":
				(*router).GET(key, httptreemux.HandlerFunc(handler))
			case "PUT":
				(*router).PUT(key, httptreemux.HandlerFunc(handler))
			case "POST":
				(*router).POST(key, httptreemux.HandlerFunc(handler))
			case "DELETE":
				(*router).DELETE(key, httptreemux.HandlerFunc(handler))
			case "OPTIONS":
				optionsHandler = handler
			}

		}
	}

	if optionsHandler != nil {
		for key, _ := range s.Routes {
			(*router).OPTIONS(key, httptreemux.HandlerFunc(optionsHandler))
		}
	}
	// generate "GET, POST, OPTIONS, HEAD, PUT" string
	for key, _ := range s.Secure.CORSMap {
		s.Secure.CORSMethods += key + ", "
	}

	// remove extra comma
	s.Secure.CORSMethods = s.Secure.CORSMethods[:len(s.Secure.CORSMethods)-2]
	s.Router = router
}