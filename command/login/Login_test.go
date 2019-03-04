package login

import (
	"github.com/ndphu/drive-manager-cli/config"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	conf := config.Config{
		BaseUrl: "https://drive-manager-api-beta.cfapps.io/api",
	}
	err := loginWithToken(&conf,
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTE3NzQwMTIsImlhdCI6MTU1MTY4NzYxMiwicHJvdmlkZXIiOiJGaXJlYmFzZSIsInR5cGUiOiJsb2dpbl90b2tlbiIsInVzZXJfZW1haWwiOiJ1c2VyMUBtYWlsaW5hdG9yLmNvbSIsInVzZXJfaWQiOiI1YzY4Yzc4OWE4OGZiNTA2ZTBkN2MwMDQifQ.9sLz3QzRsGvDq4_5p7NHrd0z0TvPCtyIf5QQcgc2LaI",
		false)
	if err != nil {
		t.Errorf("Fail to login %v\n", err)
	}
}

func TestLoginFail(t *testing.T) {
	conf := config.Config{
		BaseUrl: "https://drive-manager-api-beta.cfapps.io/api",
	}
	err := loginWithToken(&conf,
		"fake token",
		false)
	if err == nil {
		t.Errorf("Fail to login %v\n", err)
	}
}
