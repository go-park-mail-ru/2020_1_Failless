package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCaseEvents struct {
	Response   string
	StatusCode int
}

func TestFeedTEvents(t *testing.T) {
	cases := []TestCaseEvents{
		TestCaseEvents{
			Response:   `{"status": 200, "resp": {"user": 42}}`,
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
		//fmt.Println("body: ", w.Body)
		//fmt.Println("code: ", w.Code)
		//resp := w.Result()
		//body, _ := ioutil.ReadAll(resp.Body)
		//
		//bodyStr := string(body)
		//if bodyStr != item.Response {
		//	t.Errorf("[%d] wrong Response: got %+v, expected %+v",
		//		caseNum, bodyStr, item.Response)
		//}
	}
}
