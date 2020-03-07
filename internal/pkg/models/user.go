package models

import "time"

// Base user struct for storage user data
// It has a difference with the User struct into delivery package
// which contains in password filed.
type User struct {
	Uid      int    `json:"uid, omitempty"`
	Name     string `json:"name, omitempty"`
	Phone    string `json:"phone, omitempty"`
	Email    string `json:"email, omitempty"`
	Password []byte `json:"-"`
}

// Struct describes location point of user
type LocationPoint struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Accuracy  int     `json:"accuracy, omitempty"`
}

// Base profile info structure that can be used in delivery
// package for encoding/decoding json bodies
type JsonInfo struct {
	About     string        `json:"about"`
	Photos    []string      `json:"photos"`
	Rating    float32       `json:"rating"`
	Birthday  time.Time     `json:"birthday"`
	Gender    int           `json:"gender"`
	LoginDate time.Time     `json:"login_date"`
	Location  LocationPoint `json:"location"`
}
