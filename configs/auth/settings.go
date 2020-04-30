package auth

import (
	"failless/configs"
	"failless/internal/pkg/settings"
	"sync"
)

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
			Port: configs.PortAuth,
			Ip:   configs.IPAddress,
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
			MetricsHost:  configs.Prometheus,
		}
		settings.UseCaseConf = settings.GlobalConfig{
			PageLimit: 10,
			InHDD:     true,
		}
		conf.InitSecure(&settings.SecureSettings)
		conf.InitConf(&settings.UseCaseConf)
	})
	return &conf
}
