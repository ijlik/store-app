package config

import (
	"os"

	configdata "github.com/ijlik/store-app/pkg/config/data"
	configvault "github.com/ijlik/store-app/pkg/config/vault"
)

func NewConfig(
	path string,
	interval int,
) configdata.Config {
	// setup local config
	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		return getenv(path)
	}

	if interval == 0 {
		// default interval 5 second
		interval = 5
	}

	return configvault.NewVaultEnv(path, interval)
}
