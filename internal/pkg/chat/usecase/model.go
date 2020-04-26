package usecase

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/repository"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type chatUseCase struct {
	Rep chat.Repository
}

type Client struct {
	Conn    *websocket.Conn
	Request *http.Request
	Id      string
	ChatID  int64
	Uid     int64
}

func (cc *Client) Run() {
	uc := GetUseCase()
	cc.Conn.WriteJSON(network.Message{Message: cc.Id, Status: http.StatusCreated})
	lastMsgs, err := uc.GetMessagesForChat(&models.MessageRequest{
		ChatID: cc.ChatID,
		Uid:    cc.Uid,
		Limit:  20,
		Page:   0,
	})
	if err != nil {
		cc.Conn.Close()
		log.Printf("Connection %s refused: %s\n", cc.Id, err.Error())
		return
	}

	cc.Conn.WriteJSON(lastMsgs)
	for {
		var message forms.Message
		err := cc.Conn.ReadJSON(&message)
		if err != nil {
			cc.Conn.Close()
			log.Printf("Connection %s refused: %s\n", cc.Id, err.Error())
			return
		}
		user, err := security.GetUserFromCtx(cc.Request)
		if err == nil {
			message.Uid = int64(user.Uid)
		}
		message.Date = time.Now()

		message.Mid, err = uc.AddNewMessage(&message)
		if err != nil {
			log.Println(err.Error())
			err = cc.Conn.WriteJSON(
				network.Message{
					Message: err.Error(),
					Status:  http.StatusInternalServerError})
			if err != nil {
				cc.Conn.Close()
				log.Println(err.Error())
				return
			}
		}
		uc.Notify(&message)
	}
}

type Handler struct {
	Clients map[string]*Client
}

func (h *Handler) Notify(message *forms.Message) {
	var broken []string
	for _, client := range h.Clients {
		err := client.Conn.WriteJSON(message)
		if err != nil {
			client.Conn.Close()
			broken = append(broken, client.Id)
		}
	}
	// TODO: check it
	for _, detached := range broken {
		delete(h.Clients, detached)
	}
}

var MainHandler Handler

func GetUseCase() chat.UseCase {
	return &chatUseCase{
		Rep: repository.NewSqlChatRepository(db.ConnectToDB()),
	}
}

func (cc *chatUseCase) CreateDialogue(id1, id2 int) (int, error) {
	chatId, err := cc.CreateDialogue(id1, id2)

	if err != nil {
		log.Println(err)
		return -1, nil
	}

	return chatId, nil
}

func (cc *chatUseCase) IsUserHasRoom(uid int64, cid int64) (bool, error) {
	return cc.Rep.CheckRoom(cid, uid)
}

func (cc *chatUseCase) Subscribe(conn *websocket.Conn, r *http.Request) {
	id := uuid.New().String()
	cs := &Client{conn, sync.Mutex{}, r, id}
	MainHandler.Clients[id] = cs
	// TODO: implement it
	cs.Run()
}

func (cc *chatUseCase) Notify(message *forms.Message) {
}

func (cc *chatUseCase) AddNewMessage(message *forms.Message) (int64, error) {
	// check is user has this room
	has, err := cc.IsUserHasRoom(message.Uid, message.ChatID)
	if err != nil {
		log.Println("AddNewMessage: error - ", err.Error())
		return http.StatusInternalServerError, err
	}
	if !has {
		log.Println("AddNewMessage: client - room not found: ", message.Uid, message.ChatID)
		return http.StatusNotFound, nil
	}
	// insert message
	msgID, err := cc.Rep.AddMessageToChat(message, nil)
	if err != nil {
		log.Println("AddNewMessage: error while AddMessageToChat -  ", err.Error())
		return http.StatusInternalServerError, err
	}

	// TODO: check it
	message.ULocalID = msgID
	log.Println("AddNewMessage: OK")
	return http.StatusOK, nil
}

func (cc *chatUseCase) GetMessagesForChat(msgRequest *models.MessageRequest) ([]forms.Message, error) {
	has, err := cc.IsUserHasRoom(msgRequest.Uid, 0)
	if err != nil || !has {
		return nil, err
	}
	return cc.Rep.GetRoomMessages(msgRequest.Uid, msgRequest.ChatID, msgRequest.Page, msgRequest.Limit)
}

func (cc *chatUseCase) GetUserRooms(msgRequest *models.ChatRequest) ([]models.ChatMeta, error) {
	return cc.Rep.GetUserTopMessages(msgRequest.Uid, msgRequest.Page, msgRequest.Limit)
}
