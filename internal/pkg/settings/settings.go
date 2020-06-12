package settings

import (
	pb "failless/api/proto/auth"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type MapHandler struct {
	Type         string
	Handler      HandlerFunc
	AuthRequired bool
	CORS         bool
	CSRF         bool
	WS           bool
}

type GlobalSecure struct {
	CORSMethods  string
	CORSMap      map[string]struct{}
	AllowedHosts map[string]struct{}
	EnableCSRF   bool
	CSRFTokenLen int
	CSRFTokenTTL time.Duration
	MetricsHost  string
}

type GlobalConfig struct {
	PageLimit int
	InHDD     bool
}

var SecureSettings GlobalSecure
var UseCaseConf GlobalConfig
var AuthClient pb.AuthClient

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

type SMTPServer struct {
	Login	string
	Pass	string
	Host 	string
	Port 	string
	Auth	smtp.Auth
}

// Env variables which must to be set before running server
var Secrets = []string{
	"EMAIL_LOGIN",
	"EMAIL_PASSWORD",
}

var EmailServer SMTPServer

func InitSMTP() {
	EmailServer = SMTPServer{
		Host: "smtp.gmail.com",
		Port: "587",
		Login: os.Getenv(Secrets[0]),
		Pass:  os.Getenv(Secrets[1]),
	}
	EmailServer.Auth = smtp.PlainAuth("", EmailServer.Login, EmailServer.Pass, EmailServer.Host)
}
