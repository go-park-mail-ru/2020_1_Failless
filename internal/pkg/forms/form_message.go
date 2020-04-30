package forms

import (
	"html"
	"time"
)

type Message struct {
	Mid      int64     `json:"mid, omitempty"`
	Uid      int64     `json:"uid"`
	ULocalID int64     `json:"user_local_id, omitempty"`
	IsShown  bool      `json:"is_shown, omitempty"`
	ChatID   int64     `json:"chat_id, omitempty"`
	Text     string    `json:"message"`
	Date     time.Time `json:"created, omitempty"`
}

//easyjson:json
type MessageList []Message

type UserMsg struct {
	Uid    int64  `json:"uid"`
	Text   string `json:"message"`
	ChatID int64  `json:"chat_id"`
}

func (ms *Message) Validate() {
	ms.Text = html.EscapeString(ms.Text)
}
