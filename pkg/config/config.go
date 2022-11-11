package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type (
	// Specification structured configuration variables.
	Specification struct {
		App struct {
			Port string `envconfig:"APP_PORT" required:"true"`
		}

		RedisHost string `envconfig:"REDIS_HOST" required:"true"`
		RedisPort string `envconfig:"REDIS_PORT" required:"true"`
		RedisPass string `envconfig:"REDIS_PASSWORD" required:"true"`

		LogLevel string `envconfig:"LOG_LEVEL" required:"false"`

		RiderApiKey  string `envconfig:"RIDER_API_KEY" required:"true"`
		ClientApiKey string `envconfig:"CLIENT_API_KEY" required:"true"`
	}
)

// FromENV loads the environment variables to Specification.
func FromENV() (*Specification, error) {
	var config Specification
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}

func (c *Specification) AppAddress() string {
	return fmt.Sprintf("0.0.0.0:%s", c.App.Port)
}

func (c *Specification) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func (c *Specification) GetRedisPassword() string {
	return c.RedisPass
}

func (c *Specification) GetLogLevel() string {
	return c.LogLevel
}
