package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCaseEvents struct {
	Response   string
	StatusCode int
}

func TestFeedEvents(t *testing.T) {
	cases := []TestCaseEvents{
		TestCaseEvents{
			Response:   `{"status": 200, "resp": {"feed": "data"}}`,
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := "network://localhost:5000"
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		var ps map[string]string
		FeedEvents(w, req, ps)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}
	}
}
