package models

import (
	"net/http"
	"time"
)

// Gender types
const (
	Male = iota
	Female
	Other
)

type WorkMessage struct {
	Request *http.Request `json:"-"`
	Message string        `json:"message"`
	Status  int           `json:"status"`
}

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
	TagsId   []int32   `json:"tags,omitempty"`
}

//easyjson:json
type UserGeneralList []UserGeneral

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

type UserAbout struct {
	About	string	`json:"about,omitempty"`
}

type UserTags struct {
	Tags	[]int	`json:"tags, omitempty"`
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
	Tags	  []int32		`json:"tags"`
}

// For feed users
type UserRequest struct {
	Uid      int           `json:"uid, omitempty"`
	Page     int           `json:"page"`
	Limit    int           `json:"limit"`
	Query    string        `json:"query, omitempty"`
	Tags     []int         `json:"tags, omitempty"`
	Location LocationPoint `json:"location, omitempty"`
	MinAge   int           `json:"minAge"`
	MaxAge   int           `json:"maxAge"`
	Men      bool          `json:"men"`
	Women    bool          `json:"women"`
}

type ChatRoom struct {
	ChatID     int64     `json:"chat_id"`
	AdminID    int64     `json:"admin_id"`
	Created    time.Time `json:"created, omitempty"`
	UsersCount int       `json:"users_count"`
	Title      string    `json:"title, omitempty"`
	EventID    int64     `json:"event_id, omitempty"`
}

type ChatMeta struct {
	ChatID    	int64     	`json:"chat_id"`
	Name      	string    	`json:"name"`
	Avatar    	string  	`json:"avatar,omitempty"`
	Title     	string    	`json:"title,omitempty"`
	Unseen    	int       	`json:"unseen"`
	LastDate  	time.Time 	`json:"last_date"`
	LastMsg   	string    	`json:"last_msg"`
	UserCount 	int       	`json:"user_count"`
	Page      	int       	`json:"page"`
}

//easyjson:json
type ChatList []ChatMeta

type MessageRequest struct {
	ChatID int64 `json:"chat_id"`
	Uid    int64 `json:"uid"`
	Limit  int   `json:"limit"`
	Page   int   `json:"page"`
}

type ChatRequest struct {
	Uid   int64 `json:"uid"`
	Limit int   `json:"limit"`
	Page  int   `json:"page"`
}

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

//easyjson:json
type EventList []Event

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

type SmallEvent struct {
	EId    int       `json:"eid,omitempty"`
	UId    int       `json:"uid"`
	Title  string    `json:"title"`
	Descr  string    `json:"description,omitempty"`
	TagsId []int32   `json:"tags,omitempty"`
	Date   time.Time `json:"date,omitempty"`
	Photos []string  `json:"photos,omitempty"`
}

//easyjson:json
type SmallEventList []SmallEvent

type MidEvent struct {
	EId          int       `json:"eid,omitempty"`
	AdminId      int       `json:"uid"`
	Title        string    `json:"title"`
	Descr        string    `json:"description,omitempty"`
	TagsId       []int32   `json:"tags,omitempty"`
	Date         time.Time `json:"date,omitempty"`
	Photos       []string  `json:"photos,omitempty"`
	Limit        int       `json:"limit"`
	MemberAmount int       `json:"member_amount, omitempty"`
	ChatId       int       `json:"chat_id, omitempty"`
	Public       bool      `json:"public"`
	Followed     bool      `json:"followed,omitempty"`
}

//easyjson:json
type MidEventList []MidEvent

//easyjson:json
type OwnEventsList struct {
	SmallEvents SmallEventList `json:"small_events, omitempty"`
	MidEvents   MidEventList   `json:"mid_events, omitempty"`
}

type EventFollow struct {
	Uid int `json:"uid"`
	Eid int `json:"eid"`
}

//easyjson:json
type MidAndBigEventList struct {
	MidEvents MidEventList `json:"mid_events"`
	//BigEvent 	BigEventList
}

type Vote struct {
	Uid   int       `json:"uid"`
	Id    int       `json:"id"`
	Value int8      `json:"value"`
	Date  time.Time `json:"-"`
}

type Tag struct {
	Name  string `json:"name"`
	TagId int    `json:"tag_id"`
}

//easyjson:json
type TagList []Tag

//Mix up of UserGeneral, GeneralForm and Subs
type FeedPost struct {
	UserGeneral
	OnwEvents     OwnEventsList      `json:"own_events,omitempty"`
	Subscriptions MidAndBigEventList `json:"subscriptions,omitempty"`
}

//easyjson:json
type FeedResults []FeedPost

type Match struct {
	Uid     int64  `json:"uid"`
	MatchID int64  `json:"match_id"`
	Message string `json:"message"`
}
