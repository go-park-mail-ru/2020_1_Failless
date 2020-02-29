package routes


import (
	"bytes"
	"encoding/json"
	"fmt"
	// "io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"failless/server/forms"
)

type TestCaseSignUp struct {
	RegForm    forms.SignForm
	Response   string
	StatusCode int
}

func TestSignUp(t *testing.T) {
	cases := []TestCaseSignUp{
		TestCaseSignUp {
			RegForm:     forms.SignForm{
				Name:     "Sergey4Man",
				Phone:    "89999052501",
				Email:    "kerch@yndex.ru",
				Password: "full12fill",
			},
			Response:   `{"status": 200, "resp": {"user": 42}}`,
			StatusCode:  http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		result, err := json.Marshal(item.RegForm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		url := "http://localhost:5000"
		req := httptest.NewRequest("POST", url, bytes.NewBufferString(string(result)))
		w := httptest.NewRecorder()

		var ps map[string]string
		SignUp(w, req, ps)
		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}
		UserDelete(item.RegForm.Name)
		fmt.Println("body: ", w.Body)
		fmt.Println("code: ", w.Code)

	}
}


