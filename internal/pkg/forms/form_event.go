package forms

import (
	"failless/internal/pkg/models"
	"log"
	"time"
)

const (
	EventTypes       = 16 // small | middle | large
	MiddleEventLimit = 15
	TitleLenLimit    = 128
	MessageLenLimit  = 512
)

const (
	layoutISO = "2006-01-02"
	//layoutUS  = "January 1, 2020"
)

type EventForm struct {
	UId     int      `json:"uid"`
	Title   string   `json:"title"`
	Message string   `json:"description"`
	Type    int      `json:"type,omitempty"`
	Private bool     `json:"private,omitempty"`
	TagId   int      `json:"tag_id,omitempty"`
	Limit   int      `json:"limit"`
	Date    string   `json:"date"`
	Photos  []EImage `json:"photos,omitempty"`
}

func (ef *EventForm) ValidationLimits() bool {
	res := 1 <= ef.Limit && ef.Limit <= MiddleEventLimit
	log.Println(ef.Limit)
	if res {
		log.Println("limit validation ok")
	}
	return res
}

func (ef *EventForm) ValidationIDs() bool {
	res := ef.UId >= 0 && (ef.TagId >= 0 && ef.TagId <= EventTypes)
	if res {
		log.Println("validation ids ok")
	}
	return res
}

func (ef *EventForm) ValidationType() bool {
	// base check
	res := 0 <= ef.TagId && ef.TagId <= EventTypes
	if res {
		log.Println("type ok")
	}
	return res
}

func (ef *EventForm) CheckTextFields() bool {
	if ef.Title == "" {
		return false
	}

	if len(ef.Title) > TitleLenLimit {
		return false
	}

	if len(ef.Message) > MessageLenLimit {
		return false
	}
	log.Println("text fields ok")
	return true
}

func (ef *EventForm) Validate() bool {
	return ef.ValidationIDs() && ef.ValidationLimits() &&
		ef.CheckTextFields() && ef.ValidationType() // && ef.ValidationImages()
}

func (ef *EventForm) GetDBFormat(info *models.Event) {

	for _, photo := range ef.Photos {
		info.Photos = append(info.Photos, photo.ImgName)
	}

	*info = models.Event{
		AuthorId: ef.UId,
		Title:    ef.Title,
		Message:  ef.Message,
		Type:     ef.TagId,
		Limit:    ef.Limit,
	}
	if ef.Limit < 3 {
		(*info).Public = false
	} else {
		(*info).Public = true
	}

	var err error
	if ef.Date == "-" || ef.Date == "" {
		(*info).EDate = time.Now().Add(time.Hour * 8)
	} else {
		(*info).EDate, err = time.Parse(layoutISO, ef.Date)
		if err != nil {
			log.Println(ef.Date)
			log.Println(err.Error())
			return
		}
	}
	log.Println((*info).EDate)
}
