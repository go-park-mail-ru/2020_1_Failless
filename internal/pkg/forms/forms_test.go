package forms

import (
	"bytes"
	"failless/internal/pkg/models"
	"github.com/disintegration/imaging"
	"github.com/jackc/fake"
	"github.com/stretchr/testify/assert"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

const fileName = "../../../media/images/default.png"


type Faker struct {
}

func (f *Faker) GetImage() EImage {
	return EImage{
		ImgBase64: "",
		ImgName:   fake.ProductName(),
		Img:       nil,
	}
}

func (f *Faker) GetSignForm() SignForm {
	return SignForm{
		Uid:      rand.Int(),
		Name:     fake.FirstName(),
		Phone:    fake.DigitsN(8),
		Email:    fake.EmailAddress(),
		Password: fake.Password(8, 50, true, true, true),
	}
}

func (f *Faker) GetLocation() models.LocationPoint {
	return models.LocationPoint{
		Latitude:  fake.Latitute(),
		Longitude: fake.Longitude(),
		Accuracy:  rand.Int() % 100,
	}
}

func (f *Faker) GetGeneral() GeneralForm {
	return GeneralForm{
		SignForm: SignForm{},
		Tags: 	  []int32{1,2},
		Avatar:   f.GetImage(),
		Photos:   []EImage{f.GetImage()},
		Gender:   rand.Int() % 3,
		About:    fake.Paragraph(),
		Rating:   float32(rand.Int() % 5),
		Location: f.GetLocation(),
		Birthday: time.Now().Add(-time.Hour * 24 * 30 * 365 * time.Duration(rand.Int()%20)),
	}
}

func TestGeneralForm_GetDBFormat(t *testing.T) {
	faker := Faker{}
	info := models.JsonInfo{}
	user := models.User{}
	form := faker.GetGeneral()
	if err := form.GetDBFormat(&info, &user); err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, form.Birthday, info.Birthday)
	assert.Equal(t, form.Location, info.Location)
	assert.Equal(t, form.Gender, info.Gender)
	assert.Equal(t, form.Rating, info.Rating)
	assert.Equal(t, form.About, info.About)
	if len(form.Photos) != len(info.Photos) {
		t.Fail()
		return
	}

	assert.Equal(t, form.Name, user.Name)
	assert.Equal(t, form.Phone, user.Phone)
	assert.Equal(t, form.Email, user.Email)
	assert.Equal(t, form.About, info.About)
}

func TestEImage_Encode(t *testing.T) {
	eimage := EImage{}
	if err := eimage.GetImage(fileName); err != nil {
		t.Fail()
		return
	}
	err := eimage.Encode()
	if err != nil {
		t.Fail()
		return
	}
}

func TestEImage_GetImage(t *testing.T) {
	log.Println(os.Getwd())
	img, err := imaging.Open(fileName)
	if err != nil {
		t.Fail()
		return
	}
	eimage := EImage{}
	if err := eimage.GetImage(fileName); err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, eimage.Img, img)
}


func TestEImage_SaveImage(t *testing.T) {
	eimage := EImage{}
	if err := eimage.GetImage(fileName); err != nil {
		t.Fail()
		return
	}
	err := eimage.SaveImage("test")
	if err != nil {
		t.Fail()
	}
}

func TestEImage_ImageToBuffer(t *testing.T) {
	eimage := EImage{}
	if err := eimage.GetImage(fileName); err != nil {
		t.Fail()
		return
	}
	buf, err := eimage.ImageToBuffer()
	if err != nil {
		t.Fail()
		return
	}
	buf2 := new(bytes.Buffer)
	img, err := imaging.Open(fileName)
	if err != nil {
		t.Fail()
		return
	}
	err = jpeg.Encode(buf2, img, nil)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, buf, buf2)
}

func TestGeneralForm_ValidateGender(t *testing.T) {
	form := GeneralForm{Gender: models.Male}
	assert.Equal(t, true, form.ValidateGender())
	form.Gender = models.Female
	assert.Equal(t, true, form.ValidateGender())
	form.Gender = models.Other
	assert.Equal(t, true, form.ValidateGender())
	form.Gender = 4
	assert.Equal(t, false, form.ValidateGender())
}

// testing all signup forms
func TestSignForm_Validate(t *testing.T) {
	faker := Faker{}
	s := faker.GetSignForm()
	assert.Equal(t, true, s.Validate())
	s = SignForm{
		Name:     "Se",
		Phone:    "8502",
		Email:    "kerchtown!yandex.ru",
		Password: "faillesstop1indaworld",
	}
	assert.Equal(t, false, s.Validate())
	s.Email = "mrtester@test.com"
	assert.Equal(t, false, s.Validate())
	s.Phone = "9991002003"
	assert.Equal(t, false, s.Validate())
	s.Password = "failless12"
	assert.Equal(t, true, s.Validate())
}

func TestSignForm_ValidatePassword(t *testing.T) {
	s := SignForm{Password: "Fina11ySomEG00Dpass"}
	assert.Equal(t, true, s.ValidatePassword())
	s = SignForm{Password: "bad_pass"}
	assert.Equal(t, false, s.ValidatePassword())
}

func TestSignForm_ValidateEmail(t *testing.T) {
	s := SignForm{Email: fake.EmailAddress()}
	assert.Equal(t, true, s.ValidateEmail())
	s = SignForm{Email: "bad@mailru"}
	assert.Equal(t, false, s.ValidateEmail())
	s = SignForm{Email: "bad0mail.ru"}
	assert.Equal(t, false, s.ValidateEmail())
}

func TestSignForm_ValidatePhone(t *testing.T) {
	s := SignForm{Phone: "88005553535"}
	assert.Equal(t, true, s.ValidatePhone())

	s.Phone = "itsnotaphone"
	assert.Equal(t, false, s.ValidatePhone())

	s.Phone = "800"
	assert.Equal(t, false, s.ValidatePhone())

	s.Phone = "14881488148814881488"
	assert.Equal(t, false, s.ValidatePhone())

	s.Phone = "+7(495)25-25-515"
	assert.Equal(t, false, s.ValidatePhone())
	s.Phone = ""
	assert.Equal(t, false, s.ValidatePhone())
}
