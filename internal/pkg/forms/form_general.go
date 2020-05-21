package forms

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"time"
)

type GeneralForm struct {
	SignForm
	Tags     []int32         		`json:"tags,omitempty"`
	Avatar   EImage               	`json:"avatar"`
	Photos   []EImage             	`json:"photos,omitempty"`
	Gender   int                  	`json:"gender"`
	About    string               	`json:"about"`
	Rating   float32              	`json:"rating,omitempty"`
	Location models.LocationPoint 	`json:"location,omitempty"`
	Birthday time.Time            	`json:"birthday,omitempty"`
}

func (p *GeneralForm) ValidateGender() bool {
	return models.Male == p.Gender || p.Gender == models.Female || p.Gender == models.Other
}

func (p *GeneralForm) GetDBFormat(info *models.JsonInfo, user *models.User) error {
	encPass, err := security.EncryptPassword(p.Password)
	if err != nil {
		return err
	}

	var photos []string
	photos = append(photos, p.Avatar.ImgName)
	//for _, pic := range p.Photos {
	//	photos = append(photos, pic.ImgName)
	//}

	*info = models.JsonInfo{
		About:     p.About,
		Photos:    photos,
		Rating:    p.Rating,
		Birthday:  p.Birthday,
		Gender:    p.Gender,
		LoginDate: time.Time{},
		Location:  p.Location,
	}

	*user = models.User{
		Uid:      -1,
		Name:     p.Name,
		Phone:    p.Phone,
		Email:    p.Email,
		Password: encPass,
	}
	return nil
}

func (p *GeneralForm) FillProfile(row models.JsonInfo) {
	ava := ""
	if len(row.Photos) < 1 {
		ava = "default.png"
		//ava = path.Join(Media, "default.png")
	} else {
		ava = row.Photos[0]
		//ava = path.Join(Media, row.Photos[0])
	}

	p.Avatar.ImgName = ava //row.Photos[0]
	p.About = row.About
	p.Location = row.Location
	p.Gender = row.Gender
	p.Birthday = row.Birthday
	p.Rating = row.Rating
	p.Tags = row.Tags
	for _, photo := range row.Photos {
		p.Photos = append(p.Photos, EImage{ImgName: photo})
	}
}
