package forms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//testing all profile forms
func TestProfileFull(t *testing.T) {
	f := GeneralForm{
		//SignForm:  nil,
		//ImgBase64: "?",
		//ImgName:   "?",
		//Img:       nil,
		Gender:    1,
	}
	assert.Equal(t, true, f.ValidateGender())

	// needs to write GeneralForm in correct format to fully test image validation
	assert.Equal(t, false, f.ValidationImage())
}

// testing all signup forms
func TestValidateFull(t *testing.T) {
	s := SignForm{
		Name:     "Sergey",
		Phone:    "89999052502",
		Email:    "kerchtown@yandex.ru",
		Password: "faillesstop10indaworld",
	}
	assert.Equal(t, true, s.Validate())
	s = SignForm{
		Name:     "Se",
		Phone:    "8502",
		Email:    "kerchtown!yandex.ru",
		Password: "faillesstop1indaworld",
	}
	assert.Equal(t, false, s.Validate())
}


func TestValidateePassword(t *testing.T) {
	s := SignForm{
		Password: "Fina11ySomEG00Dpass",
	}
	assert.Equal(t, true, s.ValidatePassword())
	s = SignForm{
		Password: "bad_pass",
	}
	assert.Equal(t, false, s.ValidatePassword())
}


func TestValidateEmail(t *testing.T) {
	s := SignForm{
		Email: "real@mail.ru",
	}
	assert.Equal(t, true, s.ValidateEmail())
	s = SignForm{
		Email: "bad@mailru",
	}
	assert.Equal(t, false, s.ValidateEmail())
	s = SignForm{
		Email: "bad0mail.ru",
	}
	assert.Equal(t, false, s.ValidateEmail())
}

func TestValidatePhone(t *testing.T) {
	s := SignForm{
		Phone: "itsnotaphone",
	}
	assert.Equal(t, false, s.ValidatePhone())
	s = SignForm{
		Phone: "800",
	}
	assert.Equal(t, false, s.ValidatePhone())
	s = SignForm{
		Phone: "14881488148814881488",
	}
	assert.Equal(t, false, s.ValidatePhone())
	s = SignForm{
		Phone: "+7(495)25-25-515",
	}
	assert.Equal(t, false, s.ValidatePhone())
	s = SignForm{
		Phone: "88005553535",
	}
	assert.Equal(t, true, s.ValidatePhone())
}