package config

import (
	"embed"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration
type Config struct {
	Maximize     bool   `yaml:"maximize"`
	Direction    string `yaml:"direction"`
	Columns      int    `yaml:"columns"`
	OpenInNewTab bool   `yaml:"open_in_new_tab"`
}

//go:embed config.yaml
var config embed.FS

// NewConfig creates a new Config instance by parsing the yaml configuration file
func NewConfig(configPath string) (*Config, error) {
	// Check if config file exist, if not create a new one from template
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configBytes, err := config.ReadFile("config.yaml")
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve config template file: %v", err)
		}

		err = os.WriteFile("config.yaml", configBytes, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to intialize config file: %v", err)
		}
	}

	// Read from config file
	buf, err := os.ReadFile(configPath)
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
