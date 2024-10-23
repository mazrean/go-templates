package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// Addr is the address to listen on.
	Addr string
}

func NewConfig() (*Config, error) {
	viper.SetDefault("addr", ":8080")

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &c, nil
}
