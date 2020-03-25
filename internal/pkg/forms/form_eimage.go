package forms

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"failless/internal/pkg/aws"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
)

// Type for encode request for image upload
type UploadedImage struct {
	Uid      int    `json:"uid"`
	Uploaded EImage `json:"uploaded"`
}

type EImage struct {
	ImgBase64 string      `json:"img"`
	ImgName   string      `json:"path"`
	Img       image.Image `json:"-"`
}

func (pic *EImage) SaveImage(folder string) error {
	//err := imaging.Save(pic.Img, path.Join(Media, pic.ImgName))
	s3, err := aws.StartAWS()
	if err != nil {
		return err
	}
	return s3.UploadToAWS(&pic.ImgBase64, folder, pic.ImgName)
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
