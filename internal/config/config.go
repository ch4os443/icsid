package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		SSH struct {
			Port     int    `yaml:"port"`
			HostKey  string `yaml:"host_key"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"ssh"`
		Web struct {
			Port     int    `yaml:"port"`
			CertFile string `yaml:"cert_file"`
			KeyFile  string `yaml:"key_file"`
		} `yaml:"web"`
	} `yaml:"server"`
	Client struct {
		ServerAddress string `yaml:"server_address"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
	} `yaml:"client"`
	Persistence struct {
		Enabled bool   `yaml:"enabled"`
		Method  string `yaml:"method"`
		Path    string `yaml:"path"`
	} `yaml:"persistence"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
} 