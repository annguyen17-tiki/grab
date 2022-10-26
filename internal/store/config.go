package store

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server         string `required:"true" split_words:"true" default:"localhost"`
	Port           string `required:"true" split_words:"true" default:"5432"`
	DBName         string `required:"true" split_words:"true" default:"grab"`
	User           string `required:"true" split_words:"true"`
	Password       string `required:"true" split_words:"true"`
	ConnLifeTime   int    `default:"300" split_words:"true"`
	ConnTimeOut    int    `default:"30" split_words:"true"`
	MaxIdleConns   int    `default:"10" split_words:"true"`
	MaxOpenConns   int    `default:"150" split_words:"true"`
	Driver         string `default:"postgres" split_words:"true"`
	SSLMode        string `default:"disable" split_words:"true"`
	LogLevel       int    `default:"4" split_words:"true"`
	MigrateSource  string `default:"file://internal/store/migration" split_words:"true"`
	MigrateVersion uint   `required:"true" default:"20221015070000" split_words:"true"`
	MigrateTable   string `default:"migrations" split_words:"true"`
}

func (c *Config) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		c.Server, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.ConnTimeOut)
}

func loadConfig() (*Config, error) {
	var conf Config
	if err := envconfig.Process("STORE", &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
