package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCaseTags struct {
	Response   string
	StatusCode int
}

func TestFeedTags(t *testing.T) {
	cases := []TestCaseTags{
		TestCaseTags{
			Response:   `{"status": 200, "resp": {"feed": "data"}}`,
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := "http://localhost:5000"
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		var ps map[string]string
		FeedTags(w, req, ps)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}
	}
}