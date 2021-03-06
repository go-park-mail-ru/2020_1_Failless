package forms

import (
	"failless/internal/pkg/models"
	"log"
	"time"
)

type MidEventForm struct {
	AdminId    	int      	`json:"uid"`
	Title   	string   	`json:"title"`
	Descr   	string   	`json:"description,omitempty"`
	TagsId  	[]int32    	`json:"tags,omitempty"`
	Date    	time.Time	`json:"date,omitempty"`
	Photos  	[]EImage 	`json:"photos,omitempty"`
	Limit		int			`json:"limit"`
	Public		bool		`json:"public"`
}

func (mef *MidEventForm) ValidationLimits() bool {
	res := 3 <= mef.Limit && mef.Limit <= MiddleEventLimit
	if !res {
		log.Println("MIDEVENT: UserCount validation failed")
	}
	return res
}

func (mef *MidEventForm) ValidationIDs() bool {
	res := mef.AdminId > 0
	if !res {
		log.Println("MIDEVENT: AdminId validation failed")
	}
	return res
}

func (mef *MidEventForm) CheckTextFields() bool {
	if mef.Title == "" {
		log.Println("MIDEVENT: Title validation failed")
		return false
	}

	if len(mef.Title) > TitleLenLimit {
		log.Println("MIDEVENT: Title validation failed")
		return false
	}

	if len(mef.Descr) > MessageLenLimit {
		log.Println("MIDEVENT: Description validation failed")
		return false
	}
	return true
}

func (mef *MidEventForm) Validate() bool {
	return mef.ValidationIDs() && mef.ValidationLimits() &&
		mef.CheckTextFields()
}

func (mef *MidEventForm) GetDBFormat(event *models.MidEvent) {
	*event = models.MidEvent{
		AdminId: 	mef.AdminId,
		Title:    	mef.Title,
		Descr:  	mef.Descr,
		Date:		mef.Date,
		Limit:    	mef.Limit,
		Public:		mef.Public,
	}

	// Not for...range since it creates a copy of image
	for iii := 0; iii < len(mef.Photos); iii++ {
		imgName := mef.Photos[iii].ImgName
		event.Photos = append(event.Photos, imgName)
	}

	//for _, tag := range mef.TagsId {
	//	event.TagsId = append(event.TagsId, tag)
	//}
	event.TagsId = append(event.TagsId, mef.TagsId...)
}
