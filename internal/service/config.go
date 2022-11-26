package service

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BcryptCost                  int           `default:"12" split_words:"true"`
	ValidTokenHour              int           `default:"144" split_words:"true"`
	JWTSecret                   string        `default:"thisisajwtsecret" split_words:"true"`
	GeoHashPrecision            uint          `default:"6" split_words:"true"`
	MaxOffersPerBooking         int           `default:"5" split_words:"true"`
	BookingTimeout              time.Duration `default:"1m" split_words:"true"`
	CheckBookingTimeoutInterval time.Duration `default:"5s" split_words:"true"`
	RedisURL                    string        `default:"localhost" split_words:"true"`
	RedisPort                   int           `default:"6379" split_words:"true"`
	RedisPassword               string        `default:"" split_words:"true"`
	WebBaseURL                  string        `default:"http://localhost:3000" split_words:"true"`
}

func loadConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
