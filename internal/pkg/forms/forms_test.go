package forms

import (
	"failless/internal/pkg/models"
	"github.com/jackc/fake"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"testing"
	"time"
)

type Faker struct {
}

func (f *Faker) GetImage() EImage {
	return EImage{
		ImgBase64: "",
		ImgName:   fake.ProductName(),
		Img:       nil,
	}
}

func (f *Faker) GetGeneral() GeneralForm {
	return GeneralForm{
		SignForm: SignForm{},
		Events: []models.Event{
		},
		Tags: []models.Tag{

		},
		Avatar: f.GetImage(),
		Photos: []EImage{f.GetImage()},
		Gender: rand.Int() % 3,
		About:  fake.Paragraph(),
		Rating: float32(rand.Int() % 5),
		Location: models.LocationPoint{
			Latitude:  fake.Latitute(),
			Longitude: fake.Longitude(),
			Accuracy:  rand.Int() % 100,
		},
		Birthday: time.Now().Add(-time.Hour * 24 * 30 * 365 * time.Duration(rand.Int()%20)),
	}
}

func (f *Faker) GetSignForm() SignForm {
	return SignForm{
		Uid:      rand.Int(),
		Name:     fake.FirstName(),
		Phone:    fake.Phone(),
		Email:    fake.EmailAddress(),
		Password: fake.Password(4, 64, true, true, true),
	}
}

func (f *Faker) GetEventForm() EventForm {
	return EventForm{
		UId:     rand.Int(),
		Title:   fake.Title(),
		Message: fake.Paragraph(),
		Type:    rand.Int() % 16,
		Private: false,
		TagId:   rand.Int() % 16,
		Limit:   rand.Int()%15 + 2,
		Date:    "2020-10-10",
		Photos:  []EImage{f.GetImage()},
	}
}

//testing all profile forms
func TestGeneralForm_FillProfile(t *testing.T) {
	faker := Faker{}
	form := faker.GetGeneral()
	assert.Equal(t, true, form.ValidateGender())
	// needs to write GeneralForm in correct format to fully test image validation
	assert.Equal(t, false, form.ValidationImage())
}

func TestGeneralForm_GetDBFormat(t *testing.T) {
}

func TestGeneralForm_ValidateGender(t *testing.T) {
}

func TestGeneralForm_ValidationImage(t *testing.T) {
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
	assert.Equal(t, true, s.ValidatePhone())
}

func TestEventForm_CheckTextFields(t *testing.T) {
	form := EventForm{
		Title:   "",
		Message: "",
	}

	form.Title = ""
	assert.Equal(t, false, form.CheckTextFields())
	form.Title = "Ok title"
	assert.Equal(t, true, form.CheckTextFields())
	form.Title = fake.CharactersN(TitleLenLimit + 1)
	assert.Equal(t, false, form.CheckTextFields())
	form.Title = "Ok title"
	form.Message = fake.CharactersN(MessageLenLimit + 1)
	assert.Equal(t, false, form.CheckTextFields())
	form.Message = "Ok message"
	assert.Equal(t, true, form.CheckTextFields())

	faker := Faker{}
	form = faker.GetEventForm()
	assert.Equal(t, true, form.CheckTextFields())
}

func TestEventForm_ValidationType(t *testing.T) {
	form := EventForm{TagId: -1, Type: -1}
	assert.Equal(t, false, form.ValidationType())
	form.Type = EventTypes + 1
	form.TagId = EventTypes + 1
	assert.Equal(t, false, form.ValidationType())
	form.TagId = EventTypes / 2
	assert.Equal(t, true, form.ValidationType())
}

func TestEventForm_ValidationLimits(t *testing.T) {
	form := EventForm{Limit: 0}
	assert.Equal(t, false, form.ValidationLimits())
	form.Limit = MiddleEventLimit + 1
	assert.Equal(t, false, form.ValidationType())
	form.Limit = -1
	assert.Equal(t, false, form.ValidationType())
	form.Limit = MiddleEventLimit / 2
	assert.Equal(t, true, form.ValidationType())
}

func TestEventForm_ValidationIDs(t *testing.T) {
	form := EventForm{UId: -1, TagId: -1}
	assert.Equal(t, false, form.ValidationIDs())
	form.TagId = 100
	assert.Equal(t, false, form.ValidationIDs())
	form.TagId = EventTypes + 1
	form.UId = 1
	assert.Equal(t, false, form.ValidationIDs())
	form.TagId = EventTypes / 2
	assert.Equal(t, true, form.ValidationIDs())
	form.UId = -1
	assert.Equal(t, false, form.ValidationIDs())
}

func TestEventForm_GetDBFormat(t *testing.T) {
	faker := Faker{}
	form := faker.GetEventForm()
	copyForm := form
	model := models.Event{}
	form.GetDBFormat(&model)
	assert.Equal(t, true, model.AuthorId == copyForm.UId)
	assert.Equal(t, true, model.Title == copyForm.Title)
	assert.Equal(t, true, model.Message == copyForm.Message)
	assert.Equal(t, true, model.Type == copyForm.TagId)
	assert.Equal(t, true, model.Limit == copyForm.Limit)
	if copyForm.Limit < 3 {
		assert.Equal(t, false, model.Public)
	} else {
		assert.Equal(t, true, model.Public)
	}

	if copyForm.Date != "-" && copyForm.Date != "" {
		checkDate, err := time.Parse(layoutISO, copyForm.Date)
		if err != nil {
			t.Fail()
			return
		}
		assert.Equal(t, true, model.EDate == checkDate)
	}

	for i, photo := range copyForm.Photos {
		log.Println(model.Photos)
		log.Println(i)
		if len(model.Photos) == 0 || len(model.Photos) < i + 1 {
			t.Fail()
			return
		}
		assert.Equal(t, true, model.Photos[i] == photo.ImgName)
	}
}

func TestEventForm_Validate(t *testing.T) {
	faker := Faker{}
	form := faker.GetEventForm()
	assert.Equal(t, true, form.Validate())
	form.Limit = 100
	assert.Equal(t, false, form.Validate())
	form.TagId = -1
	form.Limit = 10
	assert.Equal(t, false, form.Validate())
	form.TagId = 2
	form.Type = 30
	assert.Equal(t, false, form.Validate())
	form.Type = 2
	form.Title = ""
	assert.Equal(t, false, form.Validate())
	form.Title = "Ok title"
	form.Message = fake.CharactersN(MessageLenLimit + 1)
	assert.Equal(t, false, form.Validate())
	form.Message = "Ok message"
	form.UId = -1
	assert.Equal(t, false, form.Validate())
	form.UId = 10
	assert.Equal(t, true, form.Validate())
}
