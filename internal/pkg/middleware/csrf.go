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
		log.Println("SET CSRF TOKEN", r.Method)

		if !settings.SecureSettings.EnableCSRF {
			next(w, r, ps)
			return
		}

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

// TODO: add to header CSRF-token when create request for up vote and down vote and all post|put|delete requests
// login | register pages not needed in csrf token
// and if we tape reload in the browser our SPA have to still work fine
// we get token without session, using separate method for this (get token and after this go to handler)
// we may save token to local storage in browser but not into page memory
func CheckCSRF(next settings.HandlerFunc) settings.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		log.Println("CHECK CSRF TOKEN", r.Method)

		if !settings.SecureSettings.EnableCSRF {
			next(w, r, ps)
			return
		}

		cookieToken, err := r.Cookie("csrf") // err ErrNoCookie only
		headerToken := r.Header.Get("X-CSRF-Token")
		if err != nil || headerToken != cookieToken.Value {
			log.Println("CSRF Validation Failed")
			network.GenErrorCode(w, r, "CSRF validation failed", http.StatusForbidden)
			return
		}
		next(w, r, ps)
	}
}
