package delivery

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/vote"
	"failless/internal/pkg/vote/usecase"
	"github.com/gorilla/websocket"
	json "github.com/mailru/easyjson"
	"log"
	"net/http"
)

const (
	MessageInvalidUidInBody = "uid in the body is incorrect"
)

type voteDelivery struct {
	UseCase vote.UseCase
}

func GetDelivery() vote.Delivery {
	return &voteDelivery{
		UseCase: usecase.GetUseCase(),
	}
}

func (vd *voteDelivery) VoteUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("vote for event")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var tempVote models.Vote
	err := json.UnmarshalFromReader(r.Body, &tempVote)
	if err != nil {
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	if uid != tempVote.Uid {
		network.GenErrorCode(w, r, MessageInvalidUidInBody, http.StatusBadRequest)
		return
	}

	message := vd.UseCase.VoteUser(tempVote)
	network.Jsonify(w, message, message.Status)
}

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

func (vd *voteDelivery) MatchPush(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("MatchPush: started")
	//if pusher, ok := w.(http.Pusher); ok {
	//	// Push is supported.
	//	options := &http.PushOptions{
	//		Header: http.Header{
	//			"Accept-Encoding": r.Header["Accept-Encoding"],
	//		},
	//	}
	//	if err := pusher.Push("/index.js", options); err != nil {
	//		log.Printf("ERROR - Failed to push: %v\n", err)
	//	}
	//	log.Printf("OK: Push succeed\n")
	//	network.GenErrorCode(w, r, "HTTP2 work", 200)
	//	return
	//}
	//log.Println("WARN - Pusher does not supported")
	//network.GenErrorCode(w, r, "HTTP2 failed", 200)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}


	uid := msgWithId{}
	if err = conn.ReadJSON(&uid); err != nil {
		log.Println(err)
		return
	}

	vd.UseCase.Subscribe(conn, uid.Uid)
}

