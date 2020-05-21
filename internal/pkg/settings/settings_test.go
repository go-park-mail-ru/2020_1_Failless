package settings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testServerSettings = ServerSettings{
		Port:   0,
		Ip:     "",
		Routes: nil,
		Router: nil,
		Secure: nil,
		Config: nil,
	}
	testGlobalSecure = GlobalSecure{
		CORSMethods:  "",
		CORSMap:      nil,
		AllowedHosts: nil,
		EnableCSRF:   false,
		CSRFTokenLen: 0,
		CSRFTokenTTL: 0,
		MetricsHost:  "",
	}
	testGlobalConfig = GlobalConfig{
		PageLimit: 0,
		InHDD:     true,
	}
)

func TestServerSettings_GetSettings(t *testing.T) {
	newSett := testServerSettings.GetSettings()
	assert.Equal(t, testServerSettings.Port, newSett.Port)
	assert.Equal(t, testServerSettings.Ip, newSett.Ip)
	assert.Equal(t, testServerSettings.Routes, newSett.Routes)
	assert.Equal(t, testServerSettings.Router, newSett.Router)
	assert.Equal(t, testServerSettings.Secure, newSett.Secure)
	assert.Equal(t, testServerSettings.Config, newSett.Config)
}

func TestServerSettings_InitSecure(t *testing.T) {
	tmpSett := testServerSettings
	tmpSett.InitSecure(&testGlobalSecure)
	assert.Equal(t, tmpSett.Secure, &testGlobalSecure)
}

func TestServerSettings_InitConf(t *testing.T) {
	tmpSett := testServerSettings
	tmpSett.InitConf(&testGlobalConfig)
	assert.Equal(t, tmpSett.Config, &testGlobalConfig)
}

func TestServerSettings_GetRouter(t *testing.T) {
	assert.Equal(t, testServerSettings.Router, testServerSettings.GetRouter())
}
