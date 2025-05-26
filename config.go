package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DeviceId      string `yaml:"device_id"`
	SslCaCert     string `yaml:"ssl_ca_cert"`
	SslClientCert string `yaml:"ssl_client_cert"`
	SslClientKey  string `yaml:"ssl_client_key"`
	JwtKey        string `yaml:"jwt_key"`
}

func LoadConfig(path string) (*Config, error) {
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

func MustLoadConfig(path string) *Config {
	config, err := LoadConfig(path)
	if err != nil {
		panic(err)
	}
	return config
}
