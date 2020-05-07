package images

import (
	"bytes"
	"encoding/base64"
	"failless/internal/pkg/forms"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"log"
)

const (
	Users = "users"
	App = "app"
	Events = "events"
)

// For some reason validation also uploads
// TODO: don't do it plz
func ValidateImage(image *forms.EImage, folder string) bool {
	imgBytes, err := base64.StdEncoding.DecodeString(image.ImgBase64)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	img, err := imaging.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// TODO: add normal resize
	// Resize srcImage to size = 128x128px using the Lanczos filter.
	//dstImage128 := imaging.Fill(img, 256, 256, imaging.Center, imaging.Lanczos)
	//dstImage128 = imaging.Resize(dstImage128, 256, 256, imaging.Lanczos)
	// Resize and crop the srcImage to fill the 100x100px area.
	image.Img = img // imaging.Fill(dstImage128, 100, 100, imaging.Center, imaging.Lanczos)
	image.ImgName = uuid.New().String() + ".jpg"

	err = image.SaveImage(folder)
	if err != nil {
		log.Println("Can't save image")
		log.Println(err.Error())
		return false
	}

	return true
}
