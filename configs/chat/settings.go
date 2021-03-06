package chat

import (
	"failless/configs"
	"failless/internal/pkg/chat/delivery"
	"failless/internal/pkg/router"
	"failless/internal/pkg/settings"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"google.golang.org/grpc"
	"log"
	"sync"

	pb "failless/api/proto/auth"
	mDelivery "failless/internal/pkg/metrics/delivery"
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
	"/api/chats/list": {{
		Type:         "POST",
		Handler:      delivery.GetDelivery().GetChatList,
		CORS:         true,
		AuthRequired: true,
		CSRF:         false,
		WS:           false,
	}},
	"/ws/connect": {{
		Type:         "GET",
		Handler:      delivery.GetDelivery().HandlerWS,
		CORS:         true,
		AuthRequired: false,
		CSRF:         false,
		WS:           true,
	}},
	"/api/chats/:id": {{
		Type:         "PUT",
		Handler:      delivery.GetDelivery().GetMessages,
		CORS:         true,
		AuthRequired: true,
		CSRF:         false,
		WS:           false,
	}},
	"/api/chats/:id/users": {{
		Type:         "GET",
		Handler:      delivery.GetDelivery().GetUsersForChat,
		CORS:         true,
		AuthRequired: true,
		CSRF:         false,
		WS:           false,
	}},
	"/metrics": {{
		Type:         "GET",
		Handler:      mDelivery.MetricsHandler,
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
			Port:   configs.PortChat,
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
