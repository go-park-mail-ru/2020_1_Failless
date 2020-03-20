package forms

import (
	"bytes"
	"encoding/base64"
	"failless/internal/pkg/models"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"log"
)

const (
	EventTypes       = 3 // small | middle | large
	MiddleEventLimit = 15
	TitleLenLimit    = 128
	MessageLenLimit  = 512
)

type EventForm struct {
	UId     int      `json:"uid"`
	Title   string   `json:"title"`
	Message string   `json:"description"`
	Type    int      `json:"type"`
	Private bool     `json:"private"`
	TagId   int      `json:"tag_id"`
	Limit   int      `json:"limit"`
	Photos  []EImage `json:"photos"`
}

func (ef *EventForm) ValidationImages() bool {
	for _, photo := range ef.Photos {
		imgBytes, err := base64.StdEncoding.DecodeString(photo.ImgBase64)
		if err != nil {
			log.Println(err.Error())
			return false
		}

		img, err := imaging.Decode(bytes.NewReader(imgBytes))
		if err != nil {
			log.Println(err.Error())
			return false
		}

		// Resize srcImage to size = 128x128px using the Lanczos filter.
		dstImage128 := imaging.Resize(img, 128, 128, imaging.Lanczos)
		// Resize and crop the srcImage to fill the 100x100px area.
		photo.Img = imaging.Fill(dstImage128, 100, 100, imaging.Center, imaging.Lanczos)
		photo.ImgName = uuid.New().String() + ".jpg"
		err = photo.SaveImage()
		if err != nil {
			log.Println("Can't save image")
			log.Println(err.Error())
			return false
		}
	}
	return true
}

func (ef *EventForm) ValidationLimits() bool {
	return 1 <= ef.Limit && ef.Limit < MiddleEventLimit
}

func (ef *EventForm) ValidationIDs() bool {
	return ef.UId >= 0 && ef.TagId >= 0
}

func (ef *EventForm) ValidationType() bool {
	// base check
	return 0 <= ef.Type && ef.Type <= EventTypes
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

	return true
}

func (ef *EventForm) Validate() bool {
	return ef.ValidationIDs() && ef.ValidationType() && ef.ValidationLimits() &&
		ef.CheckTextFields() && ef.ValidationImages()
}

func (ef *EventForm) GetDBFormat(info *models.Event) error {

	for _, photo := range ef.Photos {
		info.Photos = append(info.Photos, photo.ImgName)
	}

	*info = models.Event{
		AuthorId: ef.UId,
		Title:    ef.Title,
		Message:  ef.Message,
		Type:     ef.Type,
		Limit:    ef.Limit,
	}

	return nil
}
