package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	pb "failless/api/proto/auth"
	"failless/internal/pkg/security"
	"failless/internal/pkg/settings"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCaseEvents struct {
	Request    string
	Response   string
	StatusCode int
}

func TestFeedEvents(t *testing.T) {
	cases := []TestCaseEvents{
		TestCaseEvents{
			Request: "",
			Response:   `{"status": 200}`,
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := "/api/events/feed"
		req := httptest.NewRequest("GET", url, bytes.NewBuffer([]byte(item.Request)))
		w := httptest.NewRecorder()

		var ps map[string]string
		FeedEvents(w, req, ps)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}
	}
}

func TestCreateNewEvent(t *testing.T) {
	type response1 struct {
		Message string   `json:"description"`
	}
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
			Response:   `I know really nice place for go out and I like to find a company for that`,
			StatusCode: http.StatusOK,
		},
		TestCaseEvents{
			Request: `{}`,
			Response:   "",
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := "/api/event/new"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer([]byte(item.Request)))

		// Populate the request's context with our test data.
		ctx := req.Context()
		form := pb.Credentials{
			Uid:   1,
			Phone: "88005553535",
			Email: "aa@aa.aa",
			Name:  "Chort",
		}
		ctx = context.WithValue(ctx, security.CtxUserKey, form)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		var ps map[string]string
		CreateNewEvent(w, req, ps)

		res := response1{}
		json.Unmarshal([]byte(w.Body.String()), &res)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}
		if res.Message != item.Response {
			t.Errorf("[%d] wrong Response : got %s, expected %s",
				caseNum, res.Message, item.Response)
		}
	}
}

func TestGetEventsByKeyWords(t *testing.T) {
	type response1 struct {
		Message string   `json:"description"`
	}
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
			Response:   `I know really nice place for go out and I like to find a company for that`,
			StatusCode: http.StatusOK,
		},
		TestCaseEvents{
			Request: `{}`,
			Response:   `I know really nice place for go out and I like to find a company for that`,
			StatusCode: http.StatusOK,
		},
	}
	settings.UseCaseConf = settings.GlobalConfig{
		PageLimit: 1,
	}
	for caseNum, item := range cases {
		url := "/api/events/search"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer([]byte(item.Request)))

		// Populate the request's context with our test data.
		ctx := req.Context()
		form := security.UserClaims{
			Uid:   1,
			Phone: "88005553535",
			Email: "aa@aa.aa",
			Name:  "Chort",
		}
		ctx = context.WithValue(ctx, security.CtxUserKey, form)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		var ps map[string]string
		GetEventsByKeyWords(w, req, ps)

		var res []response1
		json.Unmarshal([]byte(w.Body.String()), &res)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		if res != nil {
			if res[0].Message != item.Response {
				t.Errorf("[%d] wrong Response : got %s, expected %s",
					caseNum, res[0].Message, item.Response)
			}
		}
	}
}

func TestOLDGetEventsFeed(t *testing.T) {
	type response1 struct {
		Message string   `json:"description"`
	}
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
			Response:   `I know really nice place for go out and I like to find a company for that`,
			StatusCode: http.StatusOK,
		},
		TestCaseEvents{
			Request: `{}`,
			Response:   "",
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := "/api/events/feed"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer([]byte(item.Request)))

		// Populate the request's context with our test data.
		ctx := req.Context()
		form := security.UserClaims{
			Uid:   1,
			Phone: "88005553535",
			Email: "aa@aa.aa",
			Name:  "Chort",
		}
		ctx = context.WithValue(ctx, security.CtxUserKey, form)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		var ps map[string]string
		OLDGetEventsFeed(w, req, ps)

		var res []response1
		json.Unmarshal([]byte(w.Body.String()), &res)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}
		if len(res) != 0 {
			if res[0].Message != item.Response {
				t.Errorf("[%d] wrong Response : got %s, expected %s",
					caseNum, res[0].Message, item.Response)
			}
		}
	}
}