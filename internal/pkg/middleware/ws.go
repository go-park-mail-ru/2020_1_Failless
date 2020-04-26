package middleware

import (
	"failless/internal/pkg/settings"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatHandler struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		//origin := r.Header.Get("Origin")
		//_, allowed := settings.GetAllowedOrigins()[origin]
		//return allowed
	},
}

func UpgradeToWS(next settings.HandlerFunc) settings.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		conn, _ := upgrader.Upgrade(w, r, nil)

	}
}
