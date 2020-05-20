package network

import (
	"failless/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateLogout(t *testing.T) {
	rr := httptest.NewRecorder()

	CreateLogout(rr)
	var correctCookie = http.Cookie{
		Name:     "token",
		Value:    "-",
		MaxAge:   -1,
		HttpOnly: true,
		Path: 	  "/api",
	}
	assert.Equal(t, rr.Header().Get("Set-Cookie"), correctCookie.String())
}

func TestCreateAuth(t *testing.T) {
	rr := httptest.NewRecorder()
	var user = new(models.User)

	if err := CreateAuth(rr, *user); err != nil {
		t.Fatal(err)
	}
	expires := time.Now().Add(time.Hour * 24 * 30)
	var correctCookie = http.Cookie{
		Name:     "token",
		Value:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIiwicGhvbmUiOiIiLCJlbWFpbCI6IiIsInVpZCI6MCwiZXhwIjoxNTg5MDI1Mjc3fQ.7XjgGWd7Pc1hzQKB4Oel-YDeZmbkd36tzIsR4nMzJRw",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	// TODO: cannot check for token
	assert.Equal(t, len(rr.Header().Get("Set-Cookie")), len(correctCookie.String()))
}
