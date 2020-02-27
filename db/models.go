package db

import "time"

type User struct {
	Uid      int    `json:"uid"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}

const (
	Male = iota
	Female
	Other
)

type LocationPoint struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Accuracy  int     `json:"accuracy"`
}

type UserInfo struct {
	About     string        `json:"about"`
	Photos    []string      `json:"photos"`
	Rating    float32       `json:"rating"`
	Birthday  time.Time     `json:"birthday"`
	Gender    int           `json:"gender"`
	LoginDate time.Time     `json:"login_date"`
	Location  LocationPoint `json:"location"`
}

type EType []int

type Event struct {
	EId      int       `json:"eid"`
	AuthorId int       `json:"author_id"`
	Title    string    `json:"title"`
	EDate    time.Time `json:"date"`
	Message  string    `json:"message"`
	Edited   bool      `json:"edited, omitempty"`
	Author   string    `json:"author, omitempty"`
	Type     EType     `json:"type, omitempty"`
	Limit    int       `json:"limit, omitempty"`
}

type Tag struct {
	Name  string `json:"name"`
	TagId int    `json:"tag_id"`
}
