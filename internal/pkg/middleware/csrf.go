package middleware

import (
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/settings"
	"log"
	"net/http"
)

func SetCSRF(next settings.HandlerFunc) settings.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {

		// for debug
		if !settings.SecureSettings.EnableCSRF {
			next(w, r, ps)
			return
		}

		// if it's not - set the token, if not already set
		cookie, err := r.Cookie("csrf")
		if err != nil || cookie.Value == "" {
			err = security.NewCSRFToken(w)
			if err != nil {
				log.Println("Failed to set CSRF Token")
				log.Println(err.Error())
			}
		}
		next(w, r, ps)
	}
}

func CheckCSRF(next settings.HandlerFunc) settings.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		// for debug
		if !settings.SecureSettings.EnableCSRF {
			next(w, r, ps)
			return
		}

		cookieToken, err := r.Cookie("csrf") // err ErrNoCookie only
		headerToken := r.Header.Get("X-CSRF-Token")
		if err != nil || headerToken != cookieToken.Value {
			log.Println("CSRF Validation Failed")
			network.GenErrorCode(w, r, "CSRF validation faild", http.StatusForbidden)
			return
		}
		next(w, r, ps)
	}
}
