package forms

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/disintegration/imaging"
	"github.com/go-park-mail-ru/2020_1_Failless/db"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"path"
	"time"
)

type EImage struct {
	ImgBase64 string      `json:"img"`
	ImgName   string      `json:"-"`
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
}

func (pic *EImage) GetImage(name string) (err error) {
	pic.ImgName = name
	pic.Img, err = imaging.Open(path.Join(Media, name))
	return
}

type ProfileForm struct {
	*SignForm
	Avatar   EImage           `json:"avatar"`
	Photos   []EImage         `json:"photos, omitempty"`
	Gender   int              `json:"gender"`
	About    string           `json:"about"`
	Rating   float32          `json:"rating, omitempty"`
	Location db.LocationPoint `json:"location, omitempty"`
	Birthday time.Time        `json:"birthday, omitempty"`
}

func (p *ProfileForm) ValidateGender() bool {
	return 0 <= p.Gender && p.Gender <= 2
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

	return true
}
