package delivery

import (
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

func (cd *chatDelivery) HandlerWS(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		network.GenErrorCode(w, r, network.MessageErrorWhileUpgrading, http.StatusInternalServerError)
		return
	}


	uid := msgWithId{}
	if err = conn.ReadJSON(&uid); err != nil {
		log.Println(err)
		return
	}

	cd.UseCase.Subscribe(conn, uid.Uid)
}
