package forms

import (
	"github.com/jackc/fake"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	assert.Equal(t, false, form.ValidationLimits())
	form.Limit = -1
	assert.Equal(t, false, form.ValidationLimits())
	form.Limit = MiddleEventLimit / 2
	assert.Equal(t, true, form.ValidationLimits())
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

func TestEventForm_Validate(t *testing.T) {
	form := EventForm{UId: -1, Title: "", TagId:-1, Limit: 0}
	assert.Equal(t, false, form.Validate())
	form.UId = 1
	assert.Equal(t, false, form.Validate())
	form.Limit = MiddleEventLimit / 2
	assert.Equal(t, false, form.Validate())
	form.Title = "Ok title"
	assert.Equal(t, false, form.Validate())
	form.TagId = EventTypes / 2
	assert.Equal(t, true, form.Validate())
}
