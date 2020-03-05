package server

import (
	delivery2 "failless/internal/app/auth/delivery"
	"failless/internal/app/server/delivery"
	"failless/internal/pkg/settings"
	"github.com/dimfeld/httptreemux"
	"sync"
)

var routesMap = map[string][]settings.MapHandler{
	"/api/getuser": {{
		Type:         "GET",
		Handler:      delivery.GetUserInfo,
		CORS:         false,
		AuthRequired: true,
	}},
	"/api/logout": {{
		Type:         "GET",
		Handler:      delivery2.Logout,
		CORS:         false,
		AuthRequired: true,
	}},
	"/api/signin": {{
		Type:         "POST",
		Handler:      delivery2.SignIn,
		CORS:         false,
		AuthRequired: false,
	}},
	"/api/signup": {{
		Type:         "POST",
		Handler:      delivery2.SignUp,
		CORS:         false,
		AuthRequired: false,
	}},
	"/api/events/feed": {{
		Type:         "GET",
		Handler:      delivery.FeedEvents,
		CORS:         false,
		AuthRequired: false,
	}},
	"/api/tags/feed": {{
		Type:         "GET",
		Handler:      delivery.FeedTags,
		CORS:         false,
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
			CORS:         false,
			AuthRequired: false,
		},},
}

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
			CORSMethods:  "",
			CORSMap:      map[string]struct{}{},
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
		conf.InitRouter1(httptreemux.New())
	})
	return &conf
}
