package models

import "time"

// Gender types
const (
	Male = iota
	Female
	Other
)

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

type UserGeneral struct {
	Uid      int       `json:"uid, omitempty"`
	Name     string    `json:"name, omitempty"`
	Photos   []string  `json:"photos, omitempty"`
	About    string    `json:"about, omitempty"`
	Birthday time.Time `json:"birthday, omitempty"`
	Gender   int       `json:"gender, omitempty"`
}


type DBUserGeneral struct {
	Uid      *int       `json:"uid, omitempty"`
	Name     *string    `json:"name, omitempty"`
	Photos   []string   `json:"photos, omitempty"`
	About    *string    `json:"about, omitempty"`
	Birthday *time.Time `json:"birthday, omitempty"`
	Gender   *int       `json:"gender, omitempty"`
}

func (ug *DBUserGeneral) GetUserGeneral() UserGeneral {
	user := UserGeneral{}
	user.Uid = *ug.Uid
	user.Name = *ug.Name
	if ug.Gender != nil {
		user.Gender = *ug.Gender
	}
	if ug.Birthday != nil {
		user.Birthday = *ug.Birthday
	}
	if ug.Photos != nil {
		user.Photos = ug.Photos
	}

	return user
}

// Struct describes location point of user
type LocationPoint struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
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

// For feed users
type UserRequest struct {
	Uid       int           `json:"uid, omitempty"`
	Page      int           `json:"page"`
	Limit     int           `json:"limit"`
	Query     string        `json:"query, omitempty"`
	Tags      []int         `json:"tags, omitempty"`
	Location  LocationPoint `json:"location, omitempty"`
	MinAge    int           `json:"minAge"`
	MaxAge    int           `json:"maxAge"`
	Men       bool          `json:"men"`
	Women     bool          `json:"women"`
}
