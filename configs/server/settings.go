package server

import (
	"failless/configs"
	"failless/internal/pkg/metrics/delivery"
	"failless/internal/pkg/router"
	"failless/internal/pkg/settings"
	tagDelivery "failless/internal/pkg/tag/delivery"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"google.golang.org/grpc"
	"log"
	"sync"

	pb "failless/api/proto/auth"
	eventDelivery "failless/internal/pkg/event/delivery"
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

	/***********************************************
						 Authorization
	     ***********************************************/
	"/api/srv/getuser": {{
		Type:         "GET",
		Handler:      userDelivery.GetUserInfo,
		CORS:         true,
		AuthRequired: true,
		CSRF:         false,
		WS:           false,
	}},
	"/api/srv/logout": {{
		Type:         "GET",
		Handler:      userDelivery.Logout,
		CORS:         true,
		AuthRequired: true,
		CSRF:         false,
		WS:           false,
	}},
	"/api/srv/signin": {{
		Type:         "POST",
		Handler:      userDelivery.SignIn,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
		WS:           false,
	}},
	"/api/srv/signup": {{
		Type:         "POST",
		Handler:      userDelivery.SignUp,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
		WS:           false,
	}},

	/***********************************************
	            		 Events
	***********************************************/
	"/api/srv/events/search": {{
		Type:         "POST",
		Handler:      eventDelivery.GetDelivery().GetSearchEvents,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
		WS:           false,
	}},
	"/api/srv/events/small": {
		{
			Type:         "POST",
			Handler:      eventDelivery.GetDelivery().CreateSmallEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "GET",
			Handler:      eventDelivery.GetDelivery().GetSmallEvents,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
	},
	"/api/srv/events/small/:eid": {
		{
			Type:         "PUT",
			Handler:      eventDelivery.GetDelivery().UpdateSmallEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "DELETE",
			Handler:      eventDelivery.GetDelivery().DeleteSmallEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
	},
	"/api/srv/events/mid": {{
		Type:         "POST",
		Handler:      eventDelivery.GetDelivery().CreateMiddleEvent,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/events/mid/:eid": {
		{
			Type:         "GET",
			Handler:      eventDelivery.GetDelivery().GetMiddleEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "PUT",
			Handler:      eventDelivery.GetDelivery().UpdateMiddleEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "DELETE",
			Handler:      eventDelivery.GetDelivery().DeleteMiddleEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
	},
	"/api/srv/events/mid/:eid/member": {
		{
			Type:         "POST",
			Handler:      eventDelivery.GetDelivery().JoinMiddleEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "DELETE",
			Handler:      eventDelivery.GetDelivery().LeaveMiddleEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
	},
	"/api/srv/events/big": {{
		Type:         "POST",
		Handler:      eventDelivery.GetDelivery().CreateBigEvent,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/events/big/:eid": {
		{
			Type:         "GET",
			Handler:      eventDelivery.GetDelivery().GetBigEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "PUT",
			Handler:      eventDelivery.GetDelivery().UpdateBigEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "DELETE",
			Handler:      eventDelivery.GetDelivery().DeleteBigEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
	},
	"/api/srv/events/big/:eid/visitor": {
		{
			Type:         "POST",
			Handler:      eventDelivery.GetDelivery().AddVisitorForBigEvent,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "DELETE",
			Handler:      eventDelivery.GetDelivery().RemoveVisitorForBigEvent, // TODO: create a better name
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
	},
	"/ws/match": {{
		Type:         "GET",
		Handler:      voteDelivery.GetDelivery().MatchPush,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
		WS:           true,
	}},


	/***********************************************
                 		Profile
   ***********************************************/
	"/api/srv/profile/:id/meta/about": {{
		Type:         "PUT",
		Handler:      userDelivery.UpdUserAbout,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/profile/:id/meta/tags": {{
		Type:         "PUT",
		Handler:      userDelivery.UpdUserTags,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/profile/:id/meta/photos": {{
		Type:         "PUT",
		Handler:      userDelivery.UpdUserPhotos,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/profile/:id/general": {{
		Type:         "PUT",
		Handler:      userDelivery.UpdProfileGeneral,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/profile/:id/subscriptions": {{
		Type:         "GET",
		Handler:      userDelivery.GetProfileSubscriptions,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/profile/:id": {
		{
			Type:         "PUT",
			Handler:      userDelivery.UpdProfilePage,
			CORS:         true,
			AuthRequired: true,
			CSRF:         true,
			WS:           false,
		},
		{
			Type:         "GET",
			Handler:      userDelivery.GetProfilePage,
			CORS:         true,
			AuthRequired: false,
			CSRF:         false,
			WS:           false,
		},
	},
	"/api/srv/users/:vote": {{
		Type:         "PUT",
		Handler:      voteDelivery.GetDelivery().VoteUser,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/users/feed": {{
		Type:         "POST",
		Handler:      userDelivery.GetUsersFeed,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/profile/:id/own-events": {{
		Type:         "GET",
		Handler:      userDelivery.GetSmallAndMidEventsForUser,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	"/api/srv/profile/:id/small-events": {{
		Type:         "GET",
		Handler:      userDelivery.GetSmallEventsForUser,
		CORS:         true,
		AuthRequired: true,
		CSRF:         true,
		WS:           false,
	}},
	//"/api/srv/profile/:id/mid-events": {{
	//	Type:         "GET",
	//	Handler:      userDelivery.GetOwnMidEvents,
	//	CORS:         true,
	//	AuthRequired: true,
	//	CSRF:         true,
	//	WS:           false,
	//}},

	/***********************************************
	            		 Utils
	***********************************************/
	"/api/srv/tags/feed": {{
		Type:         "GET",
		Handler:      tagDelivery.GetDelivery().FeedTags,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
		WS:           false,
	}},
	"/metrics": {{
		Type:         "GET",
		Handler:      delivery.MetricsHandler,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
		WS:           false,
	}},
	"/api": {{
		Type:         "OPTIONS",
		Handler:      router.OptionsReq,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
		WS:           false,
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
			EnableCSRF:   true,
		}
		settings.UseCaseConf = settings.GlobalConfig{
			PageLimit: 10,
			InHDD:     true,
		}
		conf.InitSecure(&settings.SecureSettings)
		conf.InitConf(&settings.UseCaseConf)
		router.InitRouter(&conf, httptreemux.New())
	})
	conn := ConnectToAuthMS(fmt.Sprintf("%s:%d", configs.AuthIP, configs.PortAuth))
	settings.AuthClient = pb.NewAuthClient(conn)
	return &conf
}
