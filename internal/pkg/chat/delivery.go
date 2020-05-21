package chat

import "net/http"

type Delivery interface {
	GetMessages(w http.ResponseWriter, r *http.Request, ps map[string]string)
	GetUsersForChat(w http.ResponseWriter, r *http.Request, ps map[string]string)
	GetChatList(w http.ResponseWriter, r *http.Request, _ map[string]string)

	// Web socket
	HandlerWS(w http.ResponseWriter, r *http.Request, _ map[string]string)
}
