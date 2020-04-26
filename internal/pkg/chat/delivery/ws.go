package delivery

import (
	"failless/internal/pkg/chat/usecase"
	"failless/internal/pkg/network"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // PepegaChamp
	},
}


func HandlerWS(w http.ResponseWriter, r *http.Request, m map[string]string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	uc := usecase.GetUseCase()

	chat := chat_logic.NewChatSocket(conn, r)
	chat.Run()
}
