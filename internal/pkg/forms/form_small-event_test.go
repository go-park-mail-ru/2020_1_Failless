package forms

import (
	"failless/internal/pkg/models"
	"github.com/jackc/fake"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSmallEventForm_ValidationIDs(t *testing.T) {
	form := SmallEventForm{Uid: -1}
	assert.Equal(t, false, form.ValidationIDs())
	form.Uid = 1
	assert.Equal(t, true, form.ValidationIDs())
}

func TestSmallEventForm_CheckTextFields(t *testing.T) {
	form := SmallEventForm{
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

func TestSmallEventForm_Validate(t *testing.T) {
	form := SmallEventForm{
		Uid: -1,
		Title:   "",
	}

	assert.Equal(t, false, form.Validate())
	form.Uid = 1
	assert.Equal(t, false, form.Validate())
	form.Title = "Ok title"
	assert.Equal(t, true, form.Validate())
}

func TestSmallEventForm_GetDBFormat(t *testing.T) {
	sef := SmallEventForm{
		Uid: 0,
		Title:   "kek",
		Descr:   "kek",
		TagsId:  nil,
		Date:    time.Time{},
		Photos:  nil,
	}
	se := models.SmallEvent{}
	sef.GetDBFormat(&se)

	assert.Equal(t, sef.Uid, se.UId)
	assert.Equal(t, sef.Title, se.Title)
	assert.Equal(t, sef.Descr, se.Descr)
	assert.Equal(t, sef.Date, se.Date)
	assert.Equal(t, len(sef.Photos), len(se.Photos))
	assert.Equal(t, len(sef.TagsId), len(se.TagsId))
}
