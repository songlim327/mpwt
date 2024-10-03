package config

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var config embed.FS

// IConfigManager defines an interface for handling config
type IConfigManager interface {
	NewConfig() (*Config, error)
	ReadConfig() (*Config, error)
	ReadConfigRaw() ([]byte, error)
	WriteConfig(config string) error
}

// Config represents the configuration
type Config struct {
	Maximize     bool   `yaml:"maximize"`
	Direction    string `yaml:"direction"`
	Columns      int    `yaml:"columns"`
	OpenInNewTab bool   `yaml:"open_in_new_tab"`
}

// ConfigManager implements the IConfigManager interface for the app config
type ConfigManager struct {
	ConfigPath string
}

// NewConfigManager creates a new ConfigManager
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{ConfigPath: configPath}
}

// NewConfig creates a new Config instance by parsing the yaml configuration file
// It will check if a config file exists, if not create a new one from template
// After the check, it will automatically call ReadConfig to read the config file
func (m *ConfigManager) NewConfig() (*Config, error) {
	if _, err := os.Stat(m.ConfigPath); os.IsNotExist(err) {
		configBytes, err := config.ReadFile("config.yaml")
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve config template file: %v", err)
		}

		err = os.WriteFile("config.yaml", configBytes, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to intialize config file: %v", err)
		}
	}

	return m.ReadConfig()
}

// ReadConfig reads config file and marshals it into Config
func (m *ConfigManager) ReadConfig() (*Config, error) {
	buf, err := m.ReadConfigRaw()
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

// WriteConfig write config string to the config file
func (m *ConfigManager) WriteConfig(config string) error {
	// Replace CRLF to LF
	buf := bytes.ReplaceAll([]byte(config), []byte("\n"), []byte("\r\n"))

	err := os.WriteFile(m.ConfigPath, buf, 0644)
	if err != nil {
		return fmt.Errorf("failed to overwrite config file: %v", err)
	}
	return nil
}

// ReadConfigRaw reads config file as raw bytes
func (m *ConfigManager) ReadConfigRaw() ([]byte, error) {
	return os.ReadFile(m.ConfigPath)
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
