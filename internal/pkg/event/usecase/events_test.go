package usecase

import (
	"bytes"
	"encoding/json"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"testing"
)

type TestCaseEvents struct {
	Request    string
}

func TestGetUseCase(t *testing.T) {

}

func TestUserUseCase_InitEventsByTime(t *testing.T) {
	var events []models.Event

	uc := GetUseCase()
	_, err := uc.InitEventsByTime(&events)

	if err != nil {
		t.Errorf("There are error in Response from repository")
	}
}

func TestUserUseCase_InitEventsByKeyWords(t *testing.T) {
	cases := []TestCaseEvents{
		TestCaseEvents{
			Request: `{
			  "uid": 1,
			  "page": 1,
			  "limit": 5,
			  "query": "I wanna go to pub",
			  "tags": [
				1
			  ],
			  "ageLimit": "20",
			  "type": 1,
			  "location": {
				"lat": 3000.2221,
				"lng": 3000.2221,
				"accurancy": 10
			  }
			}`,
		},
		TestCaseEvents{
			Request: `{}`,
		},
	}

	for caseNum, item := range cases {
		var searchRequest models.EventRequest
		decoder := json.NewDecoder(bytes.NewBuffer([]byte(item.Request)))
		_ = decoder.Decode(&searchRequest)

		var events []models.Event

		uc := GetUseCase()
		_, err := uc.InitEventsByKeyWords(&events, searchRequest.Query, searchRequest.Page)

		if caseNum == 0 && err != nil {
			t.Errorf("[%d] There are error in Response from repository",
				caseNum)
		}
		if caseNum == 1 && err == nil {
			t.Errorf("[%d] There are no error in Response from repository",
				caseNum)
		}
	}
}

func TestUserUseCase_CreateEvent(t *testing.T) {
	cases := []TestCaseEvents{
		TestCaseEvents{
			Request: `{
			  "uid": 1,
			  "title": "I wanna go to pub",
			  "date": "2020-07-07",
			  "description": "I know really nice place for go out and I like to find a company for that",
			  "type": 1,
			  "tag_id": 1,
			  "private": true,
			  "limit": 2,
			  "photos": [
				{
				  "img": "KJKJBKAKjJBKJBkjbKJBKBKJbkjbkBKbkjbbJKBKJBKb",
				  "path": "/img/defalut.png"
				}
			  ]
			}`,
		},
		TestCaseEvents{
			Request: `{}`,
		},
	}

	for caseNum, item := range cases {
		decoder := json.NewDecoder(bytes.NewBuffer([]byte(item.Request)))
		var form forms.EventForm
		_ = decoder.Decode(&form)

		uc := GetUseCase()
		_, err := uc.CreateEvent(form)

		if caseNum == 0 && err != nil {
			t.Errorf("[%d] There are error in Response from repository",
				caseNum)
		}
		if caseNum == 1 && err == nil {
			t.Errorf("[%d] There are no error in Response from repository",
				caseNum)
		}
	}
}

func TestUserUseCase_InitEventsByUserPreferences(t *testing.T) {
	cases := []TestCaseEvents{
		TestCaseEvents{
			Request: `{
			  "page": 1,
			  "limit": 5,
			  "query": "I wanna go to pub",
			  "tags": [
				1
			  ],
			  "minage": 18,
			  "maxage": 22,
			  "men": true,
			  "women": false
			}`,
		},
		TestCaseEvents{
			Request: `{}`,
		},
	}

	for caseNum, item := range cases {
		var searchRequest models.EventRequest
		decoder := json.NewDecoder(bytes.NewBuffer([]byte(item.Request)))
		_ = decoder.Decode(&searchRequest)

		var events []models.Event

		uc := GetUseCase()
		_, err := uc.InitEventsByUserPreferences(&events, &searchRequest)

		if caseNum == 0 && err != nil {
			t.Errorf("[%d] There are error in Response from repository",
				caseNum)
		}
		if caseNum == 1 && err == nil {
			t.Errorf("[%d] There are no error in Response from repository",
				caseNum)
		}
	}
}

func TestUserUseCase_TakeValidTagsOnly(t *testing.T) {

}