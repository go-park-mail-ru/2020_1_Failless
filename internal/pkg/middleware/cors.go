package middleware

import (
	"failless/internal/pkg/settings"
	"log"
	"net/http"
)

func CORS(next settings.HandlerFunc) settings.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		log.Print(r.Method, ": ", r.URL, "\n")
		origin := r.Header.Get("Origin")
		_, allowed := settings.SecureSettings.AllowedHosts[origin]
		_, allowedMethod := settings.SecureSettings.CORSMap[r.Method]
		log.Println(settings.SecureSettings.CORSMap)
		log.Println(settings.SecureSettings.AllowedHosts)
		if allowed && allowedMethod {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Set("Access-Control-Allow-Methods", settings.SecureSettings.CORSMethods)
			log.Println("This method is allowed")
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		}
		next(w, r.WithContext(r.Context()), ps)
	}
}
