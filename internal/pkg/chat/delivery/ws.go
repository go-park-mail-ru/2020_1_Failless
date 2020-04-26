package delivery

import (
	"failless/internal/pkg/chat/usecase"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandlerWS(w http.ResponseWriter, r *http.Request, m map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	uc := usecase.GetUseCase()
	uc.Subscribe(conn, r, int64(uid))
}
