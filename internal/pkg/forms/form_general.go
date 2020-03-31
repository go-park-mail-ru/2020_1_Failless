package forms

import (
	"bytes"
	"encoding/base64"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"log"
	"time"
)

type GeneralForm struct {
	SignForm
	Events   []models.Event       `json:"events"`
	Tags     []models.Tag         `json:"tags"`
	Avatar   EImage               `json:"avatar"`
	Photos   []EImage             `json:"photos, omitempty"`
	Gender   int                  `json:"gender"`
	About    string               `json:"about"`
	Rating   float32              `json:"rating, omitempty"`
	Location models.LocationPoint `json:"location, omitempty"`
	Birthday time.Time            `json:"birthday, omitempty"`
}

func (p *GeneralForm) ValidateGender() bool {
	return models.Male == p.Gender || p.Gender == models.Female || p.Gender == models.Other
}

func (p *GeneralForm) ValidationImage() bool {
	imgBytes, err := base64.StdEncoding.DecodeString(p.Avatar.ImgBase64)
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
	p.Avatar.Img = imaging.Fill(dstImage128, 100, 100, imaging.Center, imaging.Lanczos)
	p.Avatar.ImgName = uuid.New().String() + ".jpg"
	err = p.Avatar.SaveImage("users")
	if err != nil {
		log.Println("Can't save image")
		log.Println(err.Error())
		return false
	}

	p.Photos = append(p.Photos, p.Avatar)
	return true
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

func (p *GeneralForm) FillProfile(row models.JsonInfo) error {
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
	for _, photo := range row.Photos {
		p.Photos = append(p.Photos, EImage{ImgName: photo})
	}
	return nil
}