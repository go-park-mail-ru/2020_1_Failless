package server

import (
	"failless/internal/pkg/router"
	"failless/internal/pkg/settings"
	"github.com/dimfeld/httptreemux"
	"sync"

	eventDelivery "failless/internal/pkg/event/delivery"
	tagDelivery "failless/internal/pkg/tag/delivery"
	userDelivery "failless/internal/pkg/user/delivery"
)

var routesMap = map[string][]settings.MapHandler{
	"/api/getuser": {{
		Type:         "GET",
		Handler:      userDelivery.GetUserInfo,
		CORS:         true,
		AuthRequired: true,
	}},
	"/api/logout": {{
		Type:         "GET",
		Handler:      userDelivery.Logout,
		CORS:         true,
		AuthRequired: true,
	}},
	"/api/signin": {{
		Type:         "POST",
		Handler:      userDelivery.SignIn,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/signup": {{
		Type:         "POST",
		Handler:      userDelivery.SignUp,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/events/feed": {{
		Type:         "GET",
		Handler:      eventDelivery.FeedEvents,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/event/new": {{
		Type:         "POST",
		Handler:      eventDelivery.CreateNewEvent,
		CORS:         true,
		AuthRequired: true,
	}},
	"/api/search/events": {{
		Type:         "POST",
		Handler:      eventDelivery.GetEventsByKeyWords,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/tags/feed": {{
		Type:         "GET",
		Handler:      tagDelivery.FeedTags,
		CORS:         true,
		AuthRequired: false,
	}},
	"/api/profile/:id": {
		{
			Type:         "PUT",
			Handler:      userDelivery.UpdProfilePage,
			CORS:         true,
			AuthRequired: true,
		},
		{
			Type:         "GET",
			Handler:      userDelivery.GetProfilePage,
			CORS:         true,
			AuthRequired: false,
		}},
	"/api": {{
		Type:         "OPTIONS",
		Handler:      router.OptionsReq,
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
		settings.UseCaseConf = settings.GlobalConfig{
			PageLimit: 10,
		}
		conf.InitSecure(&settings.SecureSettings)
		conf.InitConf(&settings.UseCaseConf)
		router.InitRouter(&conf, httptreemux.New())
	})
	return &conf
}
