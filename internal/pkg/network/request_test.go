package network

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIdFromRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps = make(map[string]string)

	// Correct
	ps["id"] = "15"
	assert.Equal(t, 15, GetIdFromRequest(rr, req, ps))
	// Incorrect
	ps["id"] = "kek"
	assert.Equal(t, -1, GetIdFromRequest(rr, req, ps))
}

func TestGetPageFromRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps = make(map[string]string)

	// Correct
	ps["page"] = "15"
	assert.Equal(t, 15, GetPageFromRequest(rr, req, &ps))
	// Incorrect
	ps["page"] = "kek"
	assert.Equal(t, 1, GetPageFromRequest(rr, req, &ps))
}
