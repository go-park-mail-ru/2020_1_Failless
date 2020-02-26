package forms

import (
	"bytes"
	"encoding/base64"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"image"
)

type ProfileForm struct {
	*SignForm
	ImgBase64 string      `json:"img"`
	ImgName   string      `json:"-"`
	Img       image.Image `json:"-"`
	Gender    int         `json:"gender"`
}

func (p *ProfileForm) ValidateGender() bool {
	return 0 <= p.Gender && p.Gender <= 2
}

func (p *ProfileForm) ValidationImage() bool {

	imgBytes, err := base64.StdEncoding.DecodeString(p.ImgBase64)
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
	p.Img = imaging.Fill(dstImage128, 100, 100, imaging.Center, imaging.Lanczos)

	p.ImgName = uuid.New().String() + ".jpg"

	return true
}
