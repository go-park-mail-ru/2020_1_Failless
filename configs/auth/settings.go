package auth

import (
	"failless/internal/pkg/router"
	"failless/internal/pkg/settings"
	"github.com/dimfeld/httptreemux"
	"sync"

	userDelivery "failless/internal/pkg/user/delivery"
)

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
}

var doOnce sync.Once
var conf settings.ServerSettings

func GetConfig() *settings.ServerSettings {
	doOnce.Do(func() {
		conf = settings.ServerSettings{
			Port:   3002,
			Ip:     "0.0.0.0",
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
	return &conf
}
