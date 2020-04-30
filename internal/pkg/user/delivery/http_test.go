package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func signFormCheck(t *testing.T, body *bytes.Buffer, name interface{}) {
	var respForm forms.SignForm
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&respForm)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, name, respForm.Name)
	assert.Equal(t, true, respForm.Uid > 0)
}

func decodeToMsg(body *bytes.Buffer) (models.WorkMessage, error) {
	var msg models.WorkMessage
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&msg)
	if err != nil {
		return models.WorkMessage{}, err
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

	signFormCheck(t, rr.Body, mcPostBody["name"])
}

func TestSignIn(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"uid":      0,
		"name":     "mrTester",
		"phone":    "88005553535",
		"email":    "mrtester@test.com",
		"password": "qwerty12345",
	}
	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("POST", "/api/signin", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}

	rr := httptest.NewRecorder()
	var ps map[string]string
	SignIn(rr, req, ps)
	signFormCheck(t, rr.Body, mcPostBody["name"])
}

func TestLogout(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/logout", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	var ps map[string]string
	Logout(rr, req, ps)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}

	msg, err := decodeToMsg(rr.Body)
	//if err != nil {
	//	t.Fail()
	//	return
	//}
	assert.Equal(t, msg.Status, http.StatusOK)
	assert.Equal(t, msg.Message, "Successfully logout")
}

func TestGetProfilePage(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"uid":      1,
		"name":     "mrTester",
		"phone":    "88005553535",
		"email":    "mrtester@test.com",
		"password": "qwerty12345",
	}
	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("POST", "/api/profile/1", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}

	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mrtester@test.com",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)

	rr := httptest.NewRecorder()
	ps := map[string]string{}
	GetProfilePage(rr, req.WithContext(ctx), ps)
	msg, err := decodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, "Incorrect id", msg.Message)

	ps = map[string]string{"id": "1"}
	GetProfilePage(rr, req.WithContext(ctx), ps)
	decoder := json.NewDecoder(rr.Body)
	var profile forms.GeneralForm
	err = decoder.Decode(&profile)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, profile.Uid, user.Uid)
	assert.Equal(t, profile.Phone, user.Phone)
	assert.Equal(t, profile.Email, user.Email)
}

func TestUploadNewImage(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"uid":      1,
		"uploaded": "",
	}
	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("PUT", "/api/profile/1/upload", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}

	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mrtester@test.com",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)

	rr := httptest.NewRecorder()
	ps := map[string]string{"id": "1"}
	UploadNewImage(rr, req.WithContext(ctx), ps)

	msg, err := decodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, "Error within parse json", msg.Message)
}

func TestUploadNewImage2(t *testing.T) {
	eimage := map[string]interface{}{
		"img":  "long long base64 string",
		"path": "default.png",
	}
	mcPostBody := map[string]interface{}{
		"uid":      1,
		"uploaded": eimage,
	}
	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("PUT", "/api/profile/1/upload", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}

	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mrtester@test.com",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)

	rr := httptest.NewRecorder()
	ps := map[string]string{"id": "1"}
	UploadNewImage(rr, req.WithContext(ctx), ps)

	msg, err := decodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, http.StatusNotFound, msg.Status)
	assert.Equal(t, "image validation failed", msg.Message)
}

func TestUpdUserMetaData(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"uid":    1,
		"tags":   []int{1, 2, 3},
		"about":  "hello world",
		"social": []string{"vk.com/id1"},
	}
	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("PUT", "/api/profile/1/meta", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}

	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mrtester@test.com",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)

	rr := httptest.NewRecorder()
	ps := map[string]string{"id": "1"}
	UpdUserMetaData(rr, req.WithContext(ctx), ps)

	decoder := json.NewDecoder(rr.Body)
	var metaForm forms.MetaForm
	err = decoder.Decode(&metaForm)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, metaForm.Uid, user.Uid)
	assert.Equal(t, metaForm.Tags, mcPostBody["tags"])
	assert.Equal(t, metaForm.About, mcPostBody["about"])
	assert.Equal(t, metaForm.Social, mcPostBody["social"])
}

func TestUpdProfileGeneral(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"uid":      1,
		"name":     "mrTester",
		"phone":    "88005553535",
		"email":    "mrtester@test.com",
		"password": "qwerty12345",
	}

	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("PUT", "/api/profile/1/general", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}

	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mrtester@test.com",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)

	rr := httptest.NewRecorder()
	ps := map[string]string{"id": "1"}
	UpdProfileGeneral(rr, req.WithContext(ctx), ps)
	signFormCheck(t, rr.Body, user.Name)
}

func TestUpdProfilePage(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"uid":      1,
		"name":     "mrTester",
		"phone":    "88005553535",
		"email":    "mrtester@test.com",
		"password": "qwerty12345",
	}

	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("PUT", "/api/profile/1", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}

	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mrtester@test.com",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)

	rr := httptest.NewRecorder()
	ps := map[string]string{"id": "1"}
	UpdProfilePage(rr, req.WithContext(ctx), ps)

	decoder := json.NewDecoder(rr.Body)
	var generalForm forms.GeneralForm
	err = decoder.Decode(&generalForm)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, generalForm.Uid, user.Uid)
	assert.Equal(t, generalForm.Email, mcPostBody["email"])
}
