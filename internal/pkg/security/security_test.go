package security

import (
	"context"
	"failless/internal/pkg/network"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestComparePasswords(t *testing.T) {
	answer := ComparePasswords([]byte("password"), "password")
	assert.Equal(t, false, answer)
}

func TestEncryptPassword(t *testing.T) {
	_, err := EncryptPassword("password")
	assert.Equal(t, nil, err)
}

func TestCheckCredentials_Incorrect1(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	answer := CheckCredentials(rr, req)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
	assert.Equal(t, -1, answer)
}

func TestCheckCredentials_Incorrect2(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxUserKey, InvalidTestUser)

	answer := CheckCredentials(rr, req.WithContext(ctx))

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorIncorrectTokenUid, msg.Message)
	assert.Equal(t, -1, answer)
}

func TestCheckCredentials_Correct(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxUserKey, TestUser)

	answer := CheckCredentials(rr, req.WithContext(ctx))

	assert.Equal(t, TestUser.Uid, answer)
}

func TestGetUserFromCtx_Incorrect1(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	usrClaims, err := GetUserFromCtx(req)

	assert.Equal(t, UserClaims{}, usrClaims)
	assert.Equal(t, claimsNotFoundError, err)
}

func TestGetUserFromCtx_Incorrect2(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxUserKey, InvalidTestUser)

	usrClaims, err := GetUserFromCtx(req.WithContext(ctx))

	assert.Equal(t, UserClaims{}, usrClaims)
	assert.Equal(t, incorrectTokenUidError, err)
}

func TestGetUserFromCtx_Correct(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxUserKey, TestUser)

	usrClaims, err := GetUserFromCtx(req.WithContext(ctx))

	assert.Equal(t, TestUser, usrClaims)
	assert.Equal(t, nil, err)
}

func TestCompareUidsFromURLAndToken_Incorrect1(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	code := CompareUidsFromURLAndToken(rr, req, map[string]string{"id":"kek"})

	assert.Equal(t, -1, code)
}

func TestCompareUidsFromURLAndToken_Incorrect2(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxUserKey, TestUser)

	code := CompareUidsFromURLAndToken(rr, req.WithContext(ctx), map[string]string{"id":"2"})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusForbidden, msg.Status)
	assert.Equal(t, "forbidden", msg.Message)
	assert.Equal(t, -1, code)
}

func TestCompareUidsFromURLAndToken_Correct(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxUserKey, TestUser)

	code := CompareUidsFromURLAndToken(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(TestUser.Uid)})

	assert.Equal(t, TestUser.Uid, code)
}

func TestGenerateRandomBytes_Incorrect(t *testing.T) {
	bts, _ := GenerateRandomBytes(0)
	assert.Equal(t, []byte{}, bts)
}

func TestGenerateRandomBytes_Correct(t *testing.T) {
	_, err := GenerateRandomBytes(1)
	assert.Equal(t, nil, err)
}

func TestGenerateCSRFToken_Incorrect(t *testing.T) {
	str, _ := generateCSRFToken(0)
	assert.Equal(t, "", str)
}

func TestGenerateCSRFToken_Correct(t *testing.T) {
	_, err := generateCSRFToken(1)
	assert.Equal(t, nil, err)
}

func TestNewCSRFToken(t *testing.T) {
	err := NewCSRFToken(httptest.NewRecorder())
	assert.Equal(t, nil, err)
}
