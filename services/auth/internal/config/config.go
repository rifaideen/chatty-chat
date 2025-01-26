package config

import (
	"pkg/utils"
	"strconv"
)

type Config struct {
	Secret    string // jwt secret code
	ExpiresIn int    // jwt token expiry in days
}

func Load() *Config {
	// load expiry time with default value 1 day
	expiry, _ := strconv.Atoi(utils.GetEnv("JWT_EXPIRES_IN", "1"))

	return &Config{
		Secret:    utils.GetEnv("JWT_SECRET", "app-secret-code"),
		ExpiresIn: expiry,
	}
}
