package handler

import "github.com/kelseyhightower/envconfig"

var globalConfig *Config

type Config struct {
	JWTSecret string `default:"thisisajwtsecret" split_words:"true"`
}

func loadConfig() error {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return err
	}

	globalConfig = &cfg
	return nil
}

func getConfig() *Config {
	return globalConfig
}
