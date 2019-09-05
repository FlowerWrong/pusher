package config

import (
	"bufio"
	"os"

	"github.com/spf13/viper"
)

// AppEnv is app env
var AppEnv string

const (
	// DEVELOPMENT env
	DEVELOPMENT = "development"
	// TEST env
	TEST = "test"
	// PRODUCTION env
	PRODUCTION = "production"
)

// Setup ...
func Setup(file string) error {
	AppEnv = os.Getenv("APP_ENV")
	if AppEnv == "" {
		AppEnv = DEVELOPMENT
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	viper.SetConfigType("yaml")
	viper.ReadConfig(bufio.NewReader(f))

	return nil
}
