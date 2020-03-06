package server

import (
	delivery2 "failless/internal/app/auth/delivery"
	"failless/internal/app/server/delivery"
	"failless/internal/pkg/middleware"
	"failless/internal/pkg/settings"
	"github.com/dimfeld/httptreemux"
	"log"
	"sync"
)

var routesMap = map[string][]settings.MapHandler{
	"/api/getuser": {{
		Type:         "GET",
		Handler:      delivery.GetUserInfo,
		CORS:         true,
		AuthRequired: true,
	}},
	"/api/logout": {{
		Type:         "GET",
		Handler:      delivery2.Logout,
		CORS:         true,
		AuthRequired: true,
	}},
	"/api/signin": {{
		Type:         "POST",
		Handler:      delivery2.SignIn,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/signup": {{
		Type:         "POST",
		Handler:      delivery2.SignUp,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/events/feed": {{
		Type:         "GET",
		Handler:      delivery.FeedEvents,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/tags/feed": {{
		Type:         "GET",
		Handler:      delivery.FeedTags,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/profile/:id": {
		{
			Type:         "PUT",
			Handler:      delivery.UpdProfilePage,
			CORS:         true,
			AuthRequired: true,
		},
		{
			Type:         "GET",
			Handler:      delivery.GetProfilePage,
			CORS:         true,
			AuthRequired: false,
		}},
	"/api": {{
		Type:         "OPTIONS",
		Handler:      delivery.OptionsReq,
		CORS:         true,
		AuthRequired: false,
	}},
}

// Env variables which must to be set before running server
var Secrets = []string{
	"DB_NAME",
	"DB_PASSWORD",
	"DB_USER",
	//"JWT_KEY",
	//"AWS_TOKEN",
}

var doOnce sync.Once
var conf settings.ServerSettings

func GetConfig() *settings.ServerSettings {
	doOnce.Do(func() {
		conf = settings.ServerSettings{
			Port:   5000,
			Ip:     "0.0.0.0",
			Routes: routesMap,
		}
		settings.SecureSettings = settings.GlobalSecure{
			CORSMethods: "",
			CORSMap:     map[string]struct{}{},
			AllowedHosts: map[string]struct{}{
				"http://localhost":           {},
				"http://localhost:8080":      {},
				"http://localhost:5000":      {},
				"http://127.0.0.1":           {},
				"http://127.0.0.1:8080":      {},
				"http://127.0.0.1:5000":      {},
				"https://eventum.rowbot.dev": {},
			},
		}
		conf.InitSecure(&settings.SecureSettings)
		InitRouter(&conf, httptreemux.New())
	})
	return &conf
}

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
