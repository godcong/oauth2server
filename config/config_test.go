package config_test

import (
	"testing"

	"github.com/godcong/oauth2server/config"
)

func TestLoad(t *testing.T) {
	if config.DefaultConfig().GetString("database.name") == "" {
		t.Error("database.name")
	}

}
