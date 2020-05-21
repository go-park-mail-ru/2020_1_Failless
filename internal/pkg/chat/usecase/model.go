package usecase

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/repository"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	userRepository "failless/internal/pkg/user/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type chatUseCase struct {
	Rep chat.Repository
}

type Client struct {
	Conn    *websocket.Conn
	Id      string
	ChatID  []int64
	Uid     int64
}

type Handler struct {
	Clients map[string]*Client
}

func (h *Handler) Notify(message *forms.Message) {
	var broken []string
	for _, client := range h.Clients {
		for _, item := range client.ChatID {
			if item == message.ChatID {
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
	}
	// TODO: check it
	for _, detached := range broken {
		delete(h.Clients, detached)
	}
}

var MainHandler Handler

func (cc *Client) Run() {
	uc := GetUseCase()
	//msg := forms.UserMsg{}
	//err := cc.Conn.ReadJSON(&msg)
	//if err != nil {
	//	log.Println("Error with JSON unpack from SW", err)
	//	return
	//}
	//for _, room := range cc.ChatID {
	//	lastMsgs, err := uc.GetMessagesForChat(&models.MessageRequest{
	//		ChatID: room,
	//		Uid:    cc.Uid,
	//		UserCount:  10,
	//		Page:   0,
	//	})
	//	fmt.Println("uc.GetMessagesForChat(&models.MessageRequest{", lastMsgs)
	//	if err != nil {
	//		cc.Conn.Close()
	//		log.Printf("Connection %s refused: %s\n", cc.Id, err.Error())
	//		return
	//	}
	//	err = cc.Conn.WriteJSON(lastMsgs)
	//	if err != nil {
	//		cc.Conn.Close()
	//		log.Printf("Connection %s refused: %s\n", cc.Id, err.Error())
	//		return
	//	}
	//}

	for {
		message := forms.Message{}
		err := cc.Conn.ReadJSON(&message)
		//fmt.Println("cc.Conn.ReadJSON(&message)", message)

		if err != nil {
			err = cc.Conn.Close()
			if err != nil {
				log.Println(err)
			}
			log.Printf("Connection %s refused: %s\n", cc.Id, err.Error())
			return
		}

		//user, err := security.GetUserFromCtx(cc.Request)
		//if err == nil {
		//	message.Uid = int64(user.Uid)
		//}
		message.Date = time.Now()
		//fmt.Println("Before uc.AddNewMessage(&message)", message)
		code, err := uc.AddNewMessage(&message)
		if err != nil {
			log.Println(err.Error())
			err = cc.Conn.WriteJSON(
				models.WorkMessage{
					Message: err.Error(),
					Status:  code})
			if err != nil {
				err = cc.Conn.Close()
				if err != nil {
					log.Println(err)
				}
				log.Println(err.Error())
				return
			}
		}
		MainHandler.Notify(&message)
	}
}

func GetUseCase() chat.UseCase {
	return &chatUseCase{
		Rep: repository.NewSqlChatRepository(db.ConnectToDB()),
	}
}

func (cc *chatUseCase) CreateDialogue(id1, id2 int) (int, error) {
	//chatId, err := cc.CreateDialogue(id1, id2)
	//
	//if err != nil {
	//	log.Println(err)
	//	return -1, nil
	//}
	//
	//return chatId, nil
	return -1, nil
}

func (cc *chatUseCase) IsUserHasRoom(uid int64, cid int64) (bool, error) {
	return cc.Rep.CheckRoom(cid, uid)
}

func (cc *chatUseCase) Subscribe(conn *websocket.Conn, uid int64) {
	if len(MainHandler.Clients) == 0 {
		MainHandler.Clients = make(map[string]*Client)
	}

	id := uuid.New().String()
	rooms, err := cc.Rep.GetUsersRooms(uid)
	if err != nil {
		log.Println("Connection failed: ", err.Error())
		return
	}
	var roomsIDs []int64
	for _, room := range rooms {
		roomsIDs = append(roomsIDs, room.ChatID)
	}
	cs := &Client{conn, id, roomsIDs, uid}
	MainHandler.Clients[id] = cs
	cs.Run()
}

func (cc *chatUseCase) Notify(message *forms.Message) {
}

func (cc *chatUseCase) AddNewMessage(message *forms.Message) (int, error) {
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
	message.Mid = msgID
	log.Println("AddNewMessage: OK, message ID is ", message.Mid)
	return http.StatusOK, nil
}

func (cc *chatUseCase) GetMessagesForChat(msgRequest *models.MessageRequest) (forms.MessageList, error) {
	has, err := cc.IsUserHasRoom(msgRequest.Uid, msgRequest.ChatID)
	if err != nil || !has {
		return nil, err
	}
	return cc.Rep.GetRoomMessages(msgRequest.Uid, msgRequest.ChatID, msgRequest.Page, msgRequest.Limit)
}

func (cc *chatUseCase) GetUserRooms(msgRequest *models.ChatRequest) (models.ChatList, error) {
	return cc.Rep.GetUserTopMessages(msgRequest.Uid, msgRequest.Page, msgRequest.Limit)
}

func (cc *chatUseCase) GetUsersForChat(cid int64, users *models.UserGeneralList) models.WorkMessage {
	uc := userRepository.NewSqlUserRepository(db.ConnectToDB())
	return uc.GetUsersForChat(cid, users)
}
