package server

import (
	"failless/internal/app/server/delivery"
	delivery2 "failless/internal/pkg/event/delivery"
	"failless/internal/pkg/router"
	"failless/internal/pkg/settings"
	delivery4 "failless/internal/pkg/tag/delivery"
	delivery3 "failless/internal/pkg/user/delivery"
	"github.com/dimfeld/httptreemux"
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
		Handler:      delivery3.Logout,
		CORS:         true,
		AuthRequired: true,
	}},
	"/api/signin": {{
		Type:         "POST",
		Handler:      delivery3.SignIn,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/signup": {{
		Type:         "POST",
		Handler:      delivery3.SignUp,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/events/feed": {{
		Type:         "GET",
		Handler:      delivery2.FeedEvents,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/tags/feed": {{
		Type:         "GET",
		Handler:      delivery4.FeedTags,
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
		router.InitRouter(&conf, httptreemux.New())
	})
	return &conf
}
