package forms

import (
	"failless/internal/pkg/models"
	"log"
	"time"
)

type SmallEventForm struct {
	Uid	    	int      	`json:"uid"`
	Title   	string   	`json:"title"`
	Descr   	string   	`json:"description,omitempty"`
	TagsId  	[]int32    	`json:"tags,omitempty"`
	Date    	time.Time	`json:"date,omitempty"`
	Photos  	[]EImage 	`json:"photos,omitempty"`
}

func (sef *SmallEventForm) ValidationIDs() bool {
	res := sef.Uid > 0
	if !res {
		log.Println("SMALLEVENT: AdminId validation failed")
	}
	return res
}

func (sef *SmallEventForm) CheckTextFields() bool {
	if sef.Title == "" {
		log.Println("SMALLEVENT: Title validation failed")
		return false
	}

	if len(sef.Title) > TitleLenLimit {
		log.Println("SMALLEVENT: Title validation failed")
		return false
	}

	if len(sef.Descr) > MessageLenLimit {
		log.Println("SMALLEVENT: Description validation failed")
		return false
	}
	return true
}

func (sef *SmallEventForm) Validate() bool {
	return sef.ValidationIDs() && sef.CheckTextFields()
}

func (sef *SmallEventForm) GetDBFormat(event *models.SmallEvent) {
	*event = models.SmallEvent{
		UId: 		sef.Uid,
		Title:    	sef.Title,
		Descr:  	sef.Descr,
		Date:		sef.Date,
	}

	// Not for...range since it creates a copy of image
	for iii := 0; iii < len(sef.Photos); iii++ {
		imgName := sef.Photos[iii].ImgName
		event.Photos = append(event.Photos, imgName)
	}

	//for _, tag := range sef.TagsId {
	//	event.TagsId = append(event.TagsId, tag)
	//}
	event.TagsId = append(event.TagsId, sef.TagsId...)
}
