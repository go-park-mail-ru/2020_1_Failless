package forms

import (
	"html"
	"time"
)

type Message struct {
	Uid         int64     `json:"uid"`
	ULocalID    int64     `json:"u_local_id, omitempty"`
	IsShown     bool      `json:"is_shown"`
	ChatId      int64     `json:"chat_id"`
	Text        string    `json:"text"`
	Attachments []string  `json:"attachments, omitempty"`
	Date        time.Time `json:"date, omitempty"`
}

func (ms *Message) Validate() {
	ms.Text = html.EscapeString(ms.Text)
}
