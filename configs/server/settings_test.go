package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSettings(t *testing.T) {
	settings := GetSettings()
	assert.Equal(t, 5000, settings.Port)
	assert.Equal(t, "0.0.0.0", settings.Ip)

}
