package delivery

import (
	"failless/internal/pkg/network"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVoteEvent(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/event/49/like", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps map[string]string
	VoteEvent(rr, req, ps)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, msg.Status, http.StatusUnauthorized)
	assert.Equal(t, msg.Message, "auth required")
}

func TestEventFollowers(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/event/49/follow", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps map[string]string
	EventFollowers(rr, req, ps)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, msg.Status, http.StatusUnauthorized)
	assert.Equal(t, msg.Message, "auth required")
}