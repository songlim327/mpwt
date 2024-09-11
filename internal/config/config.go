package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration
type Config struct {
	Debug           bool   `yaml:"debug"`
	Maximize        bool   `yaml:"maximize"`
	Direction       string `yaml:"direction"`
	Columns         int    `yaml:"columns"`
	OpenInNewWindow bool   `yaml:"open_in_new_window"`
}

// NewConfig creates a new Config instance by parsing the yaml configuration file
func NewConfig(configFile string) (*Config, error) {
	buf, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	err = validate(c)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return c, err
}

// validate validates the configuration
func validate(c *Config) error {
	if c.Columns == 0 {
		return errors.New("columns must be specified (minimum: 1)")
	}

	if c.Direction == "" {
		return errors.New("direction must be specified (horizontal/vertical)")
	}

	return nil
}
