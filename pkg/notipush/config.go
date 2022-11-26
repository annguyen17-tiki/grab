package notipush

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	CredentialsFile   string `split_words:"true"`
	GoogleCredentials []byte `split_words:"true"`
}

func loadConfig() (*config, error) {
	var conf config
	if err := envconfig.Process("NOTIPUSH", &conf); err != nil {
		return nil, fmt.Errorf("failed to read fcm client config, err: %v", err)
	}

	return &conf, nil
}
