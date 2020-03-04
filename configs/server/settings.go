package server

import (
	delivery2 "failless/internal/app/auth/delivery"
	"failless/internal/app/server/delivery"
	"failless/internal/pkg/settings"
)

var ServerConf = settings.ServerSettings{
	Port: 5000,
	Ip:   "0.0.0.0",
	Routes: map[string][]settings.MapHandler{
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
	},
}
