package settings

import (
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type MapHandler struct {
	Type         string
	Handler      HandlerFunc
	AuthRequired bool
	CORS         bool
	CSRF         bool
}

type GlobalSecure struct {
	CORSMethods  string
	CORSMap      map[string]struct{}
	AllowedHosts map[string]struct{}
}

type GlobalConfig struct {
	PageLimit int
}

var SecureSettings GlobalSecure
var UseCaseConf GlobalConfig

type ServerSettings struct {
	Port   int
	Ip     string
	Routes map[string][]MapHandler
	Router http.Handler
	Secure *GlobalSecure
	Config *GlobalConfig
}

// Return this pointer
func (s *ServerSettings) GetSettings() ServerSettings {
	return *s
}

// Initialization of the global object with secure configurations
func (s *ServerSettings) InitSecure(secure *GlobalSecure) {
	s.Secure = secure
}

// Initialization of the global object with use case configurations
func (s *ServerSettings) InitConf(conf *GlobalConfig) {
	s.Config = conf
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

func (s *ServerSettings) GetRouter() http.Handler {
	return s.Router
}
