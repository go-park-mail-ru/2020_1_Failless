package server

import (
	delivery2 "failless/internal/app/auth/delivery"
	"failless/internal/app/server/delivery"
	"failless/internal/pkg/settings"
	"github.com/dimfeld/httptreemux"
)

var routesMap = map[string][]settings.MapHandler{
	"api/getuser": {{
		Type:    "GET",
		Handler: delivery.GetUserInfo,
	}},
	"/api/logout": {{
		Type:    "GET",
		Handler: delivery2.Logout,
	}},
	"/api/signin": {{
		Type:    "POST",
		Handler: delivery2.SignIn,
	}},
	"/api/signup": {{
		Type:    "GET",
		Handler: delivery2.SignUp,
	}},
	"/api/events/feed": {{
		Type:    "GET",
		Handler: delivery.FeedEvents,
	}},
	"/api/profile/:id": {
		{
			Type:    "POST",
			Handler: delivery.UpdProfilePage,
		},
		{
			Type:    "PUT",
			Handler: delivery.GetProfilePage,
		},},
}

var Secrets = []string{
	"DB_NAME",
	"DB_PASSWORD",
	"DB_USER",
	//"JWT_KEY",
	//"AWS_TOKEN",
}

func GetConfig() settings.ServerSettings {
	var conf = settings.ServerSettings{
		Port: 5000,
		Ip:   "0.0.0.0",
		Routes: routesMap,
	}
	conf.InitRouter1(httptreemux.New())
	return conf
}