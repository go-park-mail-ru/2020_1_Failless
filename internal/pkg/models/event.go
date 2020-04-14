package models

import "time"

type EType int

type Event struct {
	EId      int       `json:"eid"`
	AuthorId int       `json:"author_id"`
	Title    string    `json:"title"`
	EDate    time.Time `json:"date"`
	Message  string    `json:"description"`
	Edited   bool      `json:"edited, omitempty"`
	Author   string    `json:"author, omitempty"`
	Type     int       `json:"type, omitempty"`
	Limit    int       `json:"limit, omitempty"`
	Photos   []string  `json:"photos, omitempty"`
	Public   bool      `json:"public, omitempty"`
	Tag      Tag       `json:"tag, omitempty"`
}

type EventRequest struct {
	Uid       int           `json:"uid, omitempty"`
	Page      int           `json:"page"`
	Limit     int           `json:"limit"`
	UserLimit int           `json:"user_limit, omitempty"`
	Query     string        `json:"query"`
	Tags      []int         `json:"tags, omitempty"`
	Location  LocationPoint `json:"location, omitempty"`
	MinAge    int           `json:"minAge"`
	MaxAge    int           `json:"maxAge"`
	Men       bool          `json:"men"`
	Women     bool          `json:"women"`
}
