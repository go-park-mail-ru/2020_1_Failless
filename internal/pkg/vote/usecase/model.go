package usecase

//go:generate mockgen -destination=../mocks/mock_usecase.go -package=mocks failless/internal/pkg/vote UseCase

import (
	"failless/internal/pkg/chat"
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

const (
	MessageMatchHappened = "You matched with someone! Check your messages!!"
)

var (
	CorrectMessage = models.WorkMessage{
		Request: nil,
		Message: "",
		Status:  http.StatusOK,
	}
)

type voteUseCase struct {
	Rep vote.Repository
	chatRep chat.Repository
}

func GetUseCase() vote.UseCase {
	return &voteUseCase{
		Rep: repository.NewSqlVoteRepository(),
		chatRep: chatRep.NewSqlChatRepository(db.ConnectToDB()),
	}
}

func (vc *voteUseCase) VoteUser(vote models.Vote) models.WorkMessage {
	// TODO: add check is vote already be here
	vote.Value = vc.ValidateValue(vote.Value)
	err := vc.Rep.AddUserVote(vote.Uid, vote.Id, vote.Value)
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

		match, _ := vc.Rep.CheckMatching(vote.Uid, vote.Id)
		if match {
			log.Println("OK: Match happened between", vote.Uid, "and", vote.Id)
			// Create dialogue
			if _, err = vc.chatRep.InsertDialogue(
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
				count := uint8(0)
				for _, item := range clients {
					log.Println(item.Uid, matchId)
					if item.Uid == matchId {
						log.Println("Write to the channel")
						item.MessagesChannel <- models.Match{
							Uid:     matchId,
							MatchID: voteId,
							Message: "You were matched by someone",
						}
						count++
					} else if item.Uid == voteId {
						log.Println("Write to the channel")
						item.MessagesChannel <- models.Match{
							Uid:     voteId,
							MatchID: matchId,
							Message: "You've matched by someone",
						}
						count++
					}
					if count == 2 {
						break
					}

				}
			}(MainHandler.Clients, int64(vote.Id), int64(vote.Uid))

			return models.WorkMessage{
				Request: nil,
				Message: MessageMatchHappened,
				Status:  http.StatusOK,
			}
		}
	}

	return CorrectMessage
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
				err = client.Conn.Close()
				if err != nil {
					log.Println(err)
				}
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
