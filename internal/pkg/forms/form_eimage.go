package forms

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"path"
)

type EImage struct {
	ImgBase64 string      `json:"img"`
	ImgName   string      `json:"path"`
	Img       image.Image `json:"-"`
}

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

