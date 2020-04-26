package server

import (
	"failless/configs"
	"failless/internal/pkg/router"
	"failless/internal/pkg/settings"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"google.golang.org/grpc"
	"log"
	"sync"

	pb "failless/api/proto/auth"
	eventDelivery "failless/internal/pkg/event/delivery"
	tagDelivery "failless/internal/pkg/tag/delivery"
	userDelivery "failless/internal/pkg/user/delivery"
	voteDelivery "failless/internal/pkg/vote/delivery"
)

var authConn *grpc.ClientConn
var dialAuthOnce sync.Once

func ConnectToAuthMS(addr string) *grpc.ClientConn {
	dialAuthOnce.Do(
		func() {
			var err error
			authConn, err = grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Failed to dial the auth server: %v", err)
			}
		})
	return authConn
}

var routesMap = map[string][]settings.MapHandler{
	"/api/getuser": {{
		Type:         "GET",
		Handler:      userDelivery.GetUserInfo,
		CORS:         true,
		AuthRequired: true,
		CSRF:         false,
	}},
	"/api/logout": {{
		Type:         "GET",
		Handler:      userDelivery.Logout,
		CORS:         true,
		AuthRequired: true,
		CSRF:         false,
	}},
	"/api/signin": {{
		Type:         "POST",
		Handler:      userDelivery.SignIn,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
	}},
	"/api/signup": {{
		Type:         "POST",
		Handler:      userDelivery.SignUp,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
	}},
	"/api/events/feed": {
		{
			Type:         "GET",
			Handler:      eventDelivery.FeedEvents,
			CORS:         true,
			AuthRequired: false,
			CSRF:         false,
		},
		{
			Type:         "POST",
			Handler:      eventDelivery.GetEventsFeed,
			CORS:         true,
			AuthRequired: true,
			CSRF:         false,
		},
	},
	"/api/events/search": {{
		Type:         "POST",
		Handler:      eventDelivery.GetEventsByKeyWords,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
	}},
	"/api/event/new": {{
		Type:         "POST",
		Handler:      eventDelivery.CreateNewEvent,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api/event/:id/like": {{
		Type:         "POST",
		Handler:      voteDelivery.VoteEvent,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api/event/:id/dislike": {{
		Type:         "POST",
		Handler:      voteDelivery.VoteEvent,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api/event/:id/follow": {
		{
			Type:         "POST",
			Handler:      voteDelivery.FollowEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
		},
		{
			Type:         "GET",
			Handler:      voteDelivery.EventFollowers,
			CORS:         true,
			AuthRequired: true,
			CSRF:         false,
		},
	},
	"/api/tags/feed": {{
		Type:         "GET",
		Handler:      tagDelivery.FeedTags,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
	}},
	"/api/profile/:id/upload": {{
		Type:         "PUT",
		Handler:      userDelivery.UploadNewImage,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api/profile/:id/meta": {{
		Type:         "PUT",
		Handler:      userDelivery.UpdUserMetaData,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api/profile/:id/general": {{
		Type:         "PUT",
		Handler:      userDelivery.UpdProfileGeneral,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api/profile/:id/subscriptions": {{
		Type:         "GET",
		Handler:      userDelivery.GetProfileSubscriptions,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api/profile/:id": {
		{
			Type:         "PUT",
			Handler:      userDelivery.UpdProfilePage,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
		},
		{
			Type:         "GET",
			Handler:      userDelivery.GetProfilePage,
			CORS:         true,
			AuthRequired: false,
			CSRF:         false,
		},
	},
	"/api/users/:vote": {{
		Type:         "PUT",
		Handler:      voteDelivery.VoteUser,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api/users/feed": {{
		Type:         "POST",
		Handler:      userDelivery.GetUsersFeed,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
	}},
	"/api": {{
		Type:         "OPTIONS",
		Handler:      router.OptionsReq,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
	}},
}

// Env variables which must to be set before running server
var Secrets = []string{
	"POSTGRES_DB",
	"POSTGRES_PASSWORD",
	"POSTGRES_USER",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_REGION",
}

var doOnce sync.Once
var conf settings.ServerSettings

func GetConfig() *settings.ServerSettings {
	doOnce.Do(func() {
		conf = settings.ServerSettings{
			Port:   configs.PortServer,
			Ip:     configs.IPAddress,
			Routes: routesMap,
		}
		settings.SecureSettings = settings.GlobalSecure{
			CORSMethods: "",
			CORSMap:     map[string]struct{}{},
			AllowedHosts: map[string]struct{}{
				"http://localhost":           {},
				"http://localhost:8080":      {},
				"http://localhost:3000":      {},
				"http://127.0.0.1":           {},
				"http://127.0.0.1:8080":      {},
				"http://127.0.0.1:3000":      {},
				"https://eventum.rowbot.dev": {},
				"https://eventum.xyz":        {},
			},
			// referring to https://security.stackexchange.com/questions/6957/length-of-csrf-token
			// it's correct length of CSRF token for Base64 (in bytes)
			CSRFTokenLen: 20,
			CSRFTokenTTL: 1, // one hour
			EnableCSRF:   false,
		}
		settings.UseCaseConf = settings.GlobalConfig{
			PageLimit: 10,
			InHDD:     true,
		}
		conf.InitSecure(&settings.SecureSettings)
		conf.InitConf(&settings.UseCaseConf)
		router.InitRouter(&conf, httptreemux.New())
	})
	conn := ConnectToAuthMS(fmt.Sprintf("%s:%d", configs.IPAddress, configs.PortAuth))
	settings.AuthClient = pb.NewAuthClient(conn)
	return &conf
}
