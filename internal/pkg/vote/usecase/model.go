package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/vote"
	"failless/internal/pkg/vote/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"

	chatRep "failless/internal/pkg/chat/repository"
)

type voteUseCase struct {
	rep vote.Repository
}

func GetUseCase() vote.UseCase {
	return &voteUseCase{
		rep: repository.NewSqlVoteRepository(db.ConnectToDB()),
	}
}

func (vc *voteUseCase) VoteUser(vote models.Vote) models.WorkMessage {
	// TODO: add check is vote already be here
	vote.Value = vc.ValidateValue(vote.Value)
	err := vc.rep.AddUserVote(vote.Uid, vote.Id, vote.Value)
	// i think that there could be an error in one case - invalid event id
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	// Check for matching
	if vote.Value == 1 {
		log.Print("OK: check for matching\n")

		match, _ := vc.rep.CheckMatching(vote.Uid, vote.Id)
		if match {
			log.Println("OK: Match occured between", vote.Uid, "and", vote.Id)
			// Create dialogue
			cr := chatRep.NewSqlChatRepository(db.ConnectToDB())
			if _, err = cr.InsertDialogue(
				vote.Uid,
				vote.Id,
				2,
				"Чат#"+strconv.Itoa(vote.Id)); err != nil {
				log.Println(err)
				return models.WorkMessage{
					Request: nil,
					Message: err.Error(),
					Status:  http.StatusBadRequest,
				}
			}

			go func(clients map[string]*Client, voteId, matchId int64) {
				for _, item := range clients {
					log.Println(item.Uid, matchId)
					if item.Uid == matchId {
						log.Println("Write to the channel")
						item.MessagesChannel <- models.Match{
							Uid:     matchId,
							MatchID: voteId,
							Message: "You were matched by someone",
						}
						break
					}
				}
			}(MainHandler.Clients, int64(vote.Uid), int64(vote.Id))

			//MainHandler.MessagesChannel <- models.Match{
			//	Uid:     int64(vote.Id),
			//	MatchID: int64(vote.Uid),
			//	Message: "You've matched someone",
			//}

			return models.WorkMessage{
				Request: nil,
				Message: "You matched with someone! Check your messages!!",
				Status:  http.StatusOK,
			}
		}
	}

	return models.WorkMessage{
		Request: nil,
		Message: "OK",
		Status:  http.StatusOK,
	}
}

func (vc *voteUseCase) VoteEvent(vote models.Vote) models.WorkMessage {
	// TODO: add check is vote already be here
	vote.Value = vc.ValidateValue(vote.Value)
	err := vc.rep.AddEventVote(vote.Uid, vote.Id, vote.Value)
	// i think that there could be an error in one case - invalid event id
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	return models.WorkMessage{
		Request: nil,
		Message: "OK",
		Status:  http.StatusOK,
	}
}

func (vc *voteUseCase) ValidateValue(value int8) int8 {
	switch {
	case value >= 1:
		return 1
	case value <= 1:
		return -1
	}
	return 0
}

func (vc *voteUseCase) GetEventFollowers(eid int) (models.UserGeneralList, error) {
	return vc.rep.FindFollowers(eid)
}

type Client struct {
	Mut             sync.Mutex
	Conn            *websocket.Conn
	Id              string
	Uid             int64
	Cond            *sync.Cond
	MessagesChannel chan models.Match
}

func (cc *Client) Run() {
	for {
		message := <-cc.MessagesChannel
		MainHandler.Notify(&message)
	}
}

type Handler struct {
	Clients map[string]*Client
}

func (h *Handler) Notify(message *models.Match) {
	var broken []string
	for _, client := range h.Clients {
		if client.Uid == message.MatchID {
			err := client.Conn.WriteJSON(message)
			if err != nil {
				client.Conn.Close()
				broken = append(broken, client.Id)
			}
		}
	}
	// TODO: check it
	for _, detached := range broken {
		delete(h.Clients, detached)
	}
}

var MainHandler Handler

func (cc *voteUseCase) Subscribe(conn *websocket.Conn, uid int64) {
	if len(MainHandler.Clients) == 0 {
		MainHandler.Clients = make(map[string]*Client)
	}

	id := uuid.New().String()
	cs := &Client{}
	cs.Conn = conn
	cs.Id = id
	cs.Uid = uid
	cs.Cond = sync.NewCond(&cs.Mut)
	cs.MessagesChannel = make(chan models.Match)
	MainHandler.Clients[id] = cs
	cs.Run()
}
