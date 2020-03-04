package settings

import "net/http"

type handlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

// Basic router interface
type RouterInterface interface {
	NotFoundHandler(http.ResponseWriter, *http.Request)
	OptionsHandler() handlerFunc
	POST(string, handlerFunc)
	GET(string, handlerFunc)
	PUT(string, handlerFunc)
	DELETE(string, handlerFunc)
	OPTIONS(string, handlerFunc)
}

type MapHandler struct {
	Type    string
	Handler handlerFunc
}

type ServerSettings struct {
	Port   int
	Ip     string
	Routes map[string][]MapHandler
	router RouterInterface
}

// return this pointer
func (s *ServerSettings) GetSettings() ServerSettings {
	return *s
}

// Set new route
func (s *ServerSettings) SetRoute(reqType, url string, handler handlerFunc) {
	s.Routes[url] = append(s.Routes[url], MapHandler{Type: reqType, Handler: handler})
}

// Parse route map and return configured Router
func (s *ServerSettings) GetRouter() *RouterInterface {
	var optionsHandler handlerFunc = nil
	for key, list := range s.Routes {
		for _, pack := range list {
			switch pack.Type {
			case "GET":
				s.router.GET(key, pack.Handler)
			case "PUT":
				s.router.PUT(key, pack.Handler)
			case "POST":
				s.router.POST(key, pack.Handler)
			case "DELETE":
				s.router.DELETE(key, pack.Handler)
			case "OPTIONS":
				optionsHandler = pack.Handler
			}
		}
	}

	if optionsHandler != nil {
		for key, _ := range s.Routes {
			s.router.OPTIONS(key, optionsHandler)
		}
	}
	return &s.router
}
