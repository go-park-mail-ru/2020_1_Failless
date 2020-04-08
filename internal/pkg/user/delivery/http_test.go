package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//type TestCaseSignUp struct {
//	RegForm    forms.SignForm
//	StatusCode int
//}
//
//func TestSignUp(t *testing.T) {
//	url := "http://localhost:3001"
//	badForm := forms.SignForm{
//		Name:     "SergeyRof12",
//		Phone:    "28005553535",
//		Email:    "faker2@mail.ru",
//		Password: "full12fill",
//	}
//	result, _ := json.Marshal(badForm)
//	req := httptest.NewRequest("POST", url, bytes.NewBufferString(string(result)))
//	w := httptest.NewRecorder()
//	var ps map[string]string
//	SignUp(w, req, ps)
//	defer UserDelete(badForm.Email)
//
//	cases := []TestCaseSignUp{
//		TestCaseSignUp{
//			RegForm:    badForm,
//			StatusCode: http.StatusConflict, // gets ok
//		},
//		TestCaseSignUp{
//			RegForm: forms.SignForm{
//				Name:     "SergeyM1an",
//				Phone:    "89929052501",
//				Email:    "kerc2h@yndex.ru",
//				Password: "full12fill",
//			},
//			StatusCode: http.StatusOK,
//		},
//		TestCaseSignUp{
//			RegForm: forms.SignForm{
//				Name:     "F",
//				Phone:    "F",
//				Email:    "F",
//				Password: "full12fill",
//			},
//			StatusCode: http.StatusForbidden, // gets ok?
//		},
//	}
//	for caseNum, item := range cases {
//		result, err := json.Marshal(item.RegForm)
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		url := "middleware://localhost:5000"
//		req := httptest.NewRequest("POST", url, bytes.NewBufferString(string(result)))
//		w := httptest.NewRecorder()
//		SignUp(w, req, ps)
//		if w.Code != item.StatusCode {
//			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
//				caseNum, w.Code, item.StatusCode)
//		}
//		if item.RegForm != badForm {
//			UserDelete(item.RegForm.Email)
//		}
//	}
//}

func decodeToMsg(body *bytes.Buffer) (network.Message, error) {
	var msg network.Message
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&msg)
	if err != nil {
		return network.Message{}, err
	}
	return msg, nil
}

func TestGetUserInfo(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/getuser", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps map[string]string
	GetUserInfo(rr, req, ps)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
	msg, err := decodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, msg.Status, http.StatusUnauthorized)
	assert.Equal(t, msg.Message, "User is not authorised")
}

func TestSignUp(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"uid":      0,
		"name":     "mrTester",
		"phone":    "88005553535",
		"email":    "mrtester@test.com",
		"password": "qwerty12345",
	}
	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("POST", "/api/signup", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}

	rr := httptest.NewRecorder()

	var ps map[string]string
	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mail@mail.ru",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)
	SignUp(rr, req.WithContext(ctx), ps)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
	var msg network.Claims
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&msg)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, msg.Uid, user.Uid)
	assert.Equal(t, msg.Phone, user.Phone)
	assert.Equal(t, msg.Email, user.Email)


	SignUp(rr, req, ps)

	var respForm forms.SignForm
	decoder = json.NewDecoder(rr.Body)
	err = decoder.Decode(&respForm)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, respForm.Password, "")
	assert.Equal(t, respForm.Name, mcPostBody["name"])
	assert.Equal(t, true, respForm.Uid > 0)
}

