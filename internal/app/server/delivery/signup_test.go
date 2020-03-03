package delivery

import (
	"bytes"
	"encoding/json"
	"failless/server/forms"
	"fmt"
	// "io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TestCaseSignUp struct {
	RegForm    forms.SignForm
	StatusCode int
}

func TestSignUp(t *testing.T) {
	url := "network://localhost:5000"
	badForm := forms.SignForm{
		Name:     "SergeyRof12",
		Phone:    "28005553535",
		Email:    "faker2@mail.ru",
		Password: "full12fill",
	}
	result, _ := json.Marshal(badForm)
	req := httptest.NewRequest("POST", url, bytes.NewBufferString(string(result)))
	w := httptest.NewRecorder()
	var ps map[string]string
	SignUp(w, req, ps)
	defer UserDelete(badForm.Email)

	cases := []TestCaseSignUp{
		TestCaseSignUp{
			RegForm:    badForm,
			StatusCode: http.StatusConflict, // gets ok
		},
		TestCaseSignUp{
			RegForm: forms.SignForm{
				Name:     "SergeyM1an",
				Phone:    "89929052501",
				Email:    "kerc2h@yndex.ru",
				Password: "full12fill",
			},
			StatusCode: http.StatusOK,
		},
		TestCaseSignUp{
			RegForm: forms.SignForm{
				Name:     "F",
				Phone:    "F",
				Email:    "F",
				Password: "full12fill",
			},
			StatusCode: http.StatusForbidden, // gets ok?
		},
	}
	for caseNum, item := range cases {
		result, err := json.Marshal(item.RegForm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		url := "network://localhost:5000"
		req := httptest.NewRequest("POST", url, bytes.NewBufferString(string(result)))
		w := httptest.NewRecorder()
		SignUp(w, req, ps)
		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}
		if item.RegForm != badForm {
			UserDelete(item.RegForm.Email)
		}
	}
}
