package config

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

//go:embed default.yaml
var defaultConfig []byte

type Config struct {
	Debug bool
	// Addr is the address to listen on.
	Addr string
}

var configPath string

func NewConfig() (*Config, error) {
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewReader(defaultConfig)); err != nil {
		return nil, fmt.Errorf("failed to read default config: %w", err)
	}

	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("failed to merge in config: %w", err)
		}
	}

	viper.AutomaticEnv()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &c, nil
}
