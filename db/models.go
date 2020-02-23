package db

import "time"

type User struct {
	Uid      int    `json:"uid"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}

const (
	male = iota
	female
	other
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
