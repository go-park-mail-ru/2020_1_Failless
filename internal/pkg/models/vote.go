package models

import "time"

type Vote struct {
	Uid   int       `json:"uid"`
	Id    int       `json:"id"`
	Value int8      `json:"value"`
	Date  time.Time `json:"-"`
}

type Chat struct {
	ChatId    int       `json:"chat_id"`
	AdminId   int       `json:"admin_id"`
	Eid       int       `json:"eid"`
	Created   time.Time `json:"created, omitempty"`
	UserCount int       `json:"user_count, omitempty"`
	Title     string    `json:"title, omitempty"`
}
