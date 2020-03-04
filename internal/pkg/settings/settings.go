package settings

import (
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type MapHandler struct {
	Type    string
	Handler HandlerFunc
}

type ServerSettings struct {
	Port   int
	Ip     string
	Routes map[string][]MapHandler
	Router http.Handler
}

// return this pointer
func (s *ServerSettings) GetSettings() ServerSettings {
	return *s
}

// Set new route
func (s *ServerSettings) SetRoute(reqType, url string, handler HandlerFunc) {
	s.Routes[url] = append(s.Routes[url], MapHandler{Type: reqType, Handler: handler})
}

func (s *ServerSettings) SetRouter(handler http.Handler) {
	s.Router = handler
}

// Basic Router interface
type RouterInterface interface {
	http.Handler
	POST(path string, handler HandlerFunc)
	GET(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	OPTIONS(path string, handler HandlerFunc)
}

// Parse route map and return configured Router
func (s *ServerSettings) InitRouter1(router *httptreemux.TreeMux) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Error was occurred", r)
		}
	}()
	var optionsHandler HandlerFunc = nil
	for key, list := range s.Routes {
		log.Println(key)
		for _, pack := range list {
			switch pack.Type {
			case "GET":
				(*router).GET(key, httptreemux.HandlerFunc(pack.Handler))
			case "PUT":
				(*router).PUT(key, httptreemux.HandlerFunc(pack.Handler))
			case "POST":
				(*router).POST(key, httptreemux.HandlerFunc(pack.Handler))
			case "DELETE":
				(*router).DELETE(key, httptreemux.HandlerFunc(pack.Handler))
			case "OPTIONS":
				optionsHandler = pack.Handler
			}
		}
	}

	if optionsHandler != nil {
		for key, _ := range s.Routes {
			(*router).OPTIONS(key, httptreemux.HandlerFunc(optionsHandler))
		}
	}
	s.Router = router
}

func (s *ServerSettings) GetRouter() http.Handler {
	return s.Router
}
