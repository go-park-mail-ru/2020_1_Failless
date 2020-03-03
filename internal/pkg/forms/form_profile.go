package forms

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"failless/internal/pkg/db"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"path"
	"time"
)

type EImage struct {
	ImgBase64 string      `json:"img"`
	ImgName   string      `json:"path"`
	Img       image.Image `json:"-"`
}

const (
	Media = "media/images/"
)

func (pic *EImage) SaveImage() error {
	err := imaging.Save(pic.Img, path.Join(Media, pic.ImgName))
	return err
}

func (pic *EImage) Encode() error {
	buf := make([]byte, 128*128*3)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := jpeg.Encode(w, pic.Img, &jpeg.Options{})
	if err != nil {
		return err
	}
	pic.ImgBase64 = base64.StdEncoding.EncodeToString(buf)
	return nil
}

func (pic *EImage) GetImage(name string) (err error) {
	pic.ImgName = name
	pic.Img, err = imaging.Open(name)
	return
}

type ProfileForm struct {
	SignForm
	Avatar   EImage           `json:"avatar"`
	Photos   []EImage         `json:"photos, omitempty"`
	Gender   int              `json:"gender"`
	About    string           `json:"about"`
	Rating   float32          `json:"rating, omitempty"`
	Location db.LocationPoint `json:"location, omitempty"`
	Birthday time.Time        `json:"birthday, omitempty"`
}

func (p *ProfileForm) ValidateGender() bool {
	return db.Male == p.Gender || p.Gender == db.Female || p.Gender == db.Other
}

func (p *ProfileForm) ValidationImage() bool {

	imgBytes, err := base64.StdEncoding.DecodeString(p.Avatar.ImgBase64)
	if err != nil {
		return false
	}

	img, err := imaging.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return false
	}
	// Resize srcImage to size = 128x128px using the Lanczos filter.
	dstImage128 := imaging.Resize(img, 128, 128, imaging.Lanczos)

	// Resize and crop the srcImage to fill the 100x100px area.
	p.Avatar.Img = imaging.Fill(dstImage128, 100, 100, imaging.Center, imaging.Lanczos)

	p.Avatar.ImgName = uuid.New().String() + ".jpg"
	err = p.Avatar.SaveImage()
	if err != nil {
		return false
	}

	p.Photos = append(p.Photos, p.Avatar)
	return true
}

func (p *ProfileForm) GetDBFormat(info *db.UserInfo, user *db.User) error {
	encPass, err := EncryptPassword(p.Password)
	if err != nil {
		return err
	}

	var photos []string
	photos = append(photos, p.Avatar.ImgName)
	//for _, pic := range p.Photos {
	//	photos = append(photos, pic.ImgName)
	//}

	*info = db.UserInfo{
		About:     p.About,
		Photos:    photos,
		Rating:    p.Rating,
		Birthday:  p.Birthday,
		Gender:    p.Gender,
		LoginDate: time.Time{},
		Location:  p.Location,
	}
	*user = db.User{
		Uid:      -1,
		Name:     p.Name,
		Phone:    p.Phone,
		Email:    p.Email,
		Password: encPass,
	}
	return nil
}

func (p *ProfileForm) FillProfile(row db.UserInfo) error {
	// todo: take pictures from media
	ava := ""
	if len(row.Photos) < 1 {
		ava = "default.png"
		//ava = path.Join(Media, "default.png")
	} else {
		ava = row.Photos[0]
		//ava = path.Join(Media, row.Photos[0])
	}
	//if err := eimage.GetImage(ava); err != nil {
	//	return err
	//}
	p.Avatar.ImgName = ava //row.Photos[0]
	p.About = row.About
	p.Location = row.Location
	p.Gender = row.Gender
	return nil
}
