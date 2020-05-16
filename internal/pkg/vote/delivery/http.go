package delivery

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/vote/usecase"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	json "github.com/mailru/easyjson"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////// MANAGE ///////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////
func VoteUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("vote for event")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var vote models.Vote
	err := json.UnmarshalFromReader(r.Body, &vote)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	if uid != vote.Uid {
		network.GenErrorCode(w, r, "uid in the body is incorrect", http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	message := uc.VoteUser(vote)
	network.Jsonify(w, message, message.Status)
}

func VoteEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("vote for event")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var vote models.Vote
	err := json.UnmarshalFromReader(r.Body, &vote)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	if uid != vote.Uid {
		network.GenErrorCode(w, r, "uid in the body is incorrect", http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	message := uc.VoteEvent(vote)
	network.Jsonify(w, message, message.Status)
}

func EventFollowers(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("follow for event")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	id := network.GetIdFromRequest(w, r, ps)
	if id < 0 {
		network.GenErrorCode(w, r, "url id is incorrect", http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	followers, err := uc.GetEventFollowers(int(id))
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	network.Jsonify(w, followers, http.StatusOK)
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

func MatchPush(w http.ResponseWriter, r *http.Request, ps map[string]string) {
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

	uc := usecase.GetUseCase()
	uc.Subscribe(conn, uid.Uid)
}
