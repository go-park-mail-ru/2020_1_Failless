package forms

import (
	"failless/internal/pkg/models"
	"github.com/jackc/fake"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMidEventForm_ValidationLimits(t *testing.T) {
	form := MidEventForm{Limit: 0}
	assert.Equal(t, false, form.ValidationLimits())
	form.Limit = MiddleEventLimit + 1
	assert.Equal(t, false, form.ValidationLimits())
	form.Limit = -1
	assert.Equal(t, false, form.ValidationLimits())
	form.Limit = MiddleEventLimit / 2
	assert.Equal(t, true, form.ValidationLimits())
}

func TestMidEventForm_ValidationIDs(t *testing.T) {
	form := MidEventForm{AdminId: -1}
	assert.Equal(t, false, form.ValidationIDs())
	form.AdminId = 1
	assert.Equal(t, true, form.ValidationIDs())
}

func TestMidEventForm_CheckTextFields(t *testing.T) {
	form := MidEventForm{
		Title:   "",
		Descr: "",
	}

	form.Title = ""
	assert.Equal(t, false, form.CheckTextFields())
	form.Title = "Ok title"
	assert.Equal(t, true, form.CheckTextFields())
	form.Title = fake.CharactersN(TitleLenLimit + 1)
	assert.Equal(t, false, form.CheckTextFields())
	form.Title = "Ok title"
	form.Descr = fake.CharactersN(MessageLenLimit + 1)
	assert.Equal(t, false, form.CheckTextFields())
	form.Descr = "Ok message"
	assert.Equal(t, true, form.CheckTextFields())
}

func TestMidEventForm_Validate(t *testing.T) {
	form := MidEventForm{
		AdminId: -1,
		Title:   "",
		Descr: "",
		Limit: 0,
	}

	assert.Equal(t, false, form.Validate())
	form.AdminId = 1
	assert.Equal(t, false, form.Validate())
	form.Limit = MiddleEventLimit / 2
	assert.Equal(t, false, form.Validate())
	form.Title = "Ok title"
	assert.Equal(t, true, form.Validate())
}

func TestMidEventForm_GetDBFormat(t *testing.T) {
	mef := MidEventForm{
		AdminId: 0,
		Title:   "kek",
		Descr:   "kek",
		TagsId:  nil,
		Date:    time.Time{},
		Photos:  nil,
		Limit:   0,
		Public:  false,
	}
	me := models.MidEvent{}
	mef.GetDBFormat(&me)

	assert.Equal(t, mef.AdminId, me.AdminId)
	assert.Equal(t, mef.Title, me.Title)
	assert.Equal(t, mef.Descr, me.Descr)
	assert.Equal(t, mef.Date, me.Date)
	assert.Equal(t, mef.Limit, me.Limit)
	assert.Equal(t, mef.Public, me.Public)
	assert.Equal(t, len(mef.Photos), len(me.Photos))
	assert.Equal(t, len(mef.TagsId), len(me.TagsId))
}
