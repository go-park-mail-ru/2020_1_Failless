package delivery

import (
	"failless/internal/pkg/chat/usecase"
	"failless/internal/pkg/network"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type msgWithId struct {
	Uid 	int64
}

func HandlerWS(w http.ResponseWriter, r *http.Request, m map[string]string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}


	uid := msgWithId{}
	if err = conn.ReadJSON(&uid); err != nil {
		log.Println(err)
		return
	}

	uc := usecase.GetUseCase()
	uc.Subscribe(conn, r, uid.Uid)
}
