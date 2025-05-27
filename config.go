package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DeviceId      string `yaml:"device_id"`
	DeviceName    string `yaml:"device_name"`
	SslCaCert     string `yaml:"ssl_ca_cert"`
	SslClientCert string `yaml:"ssl_client_cert"`
	SslClientKey  string `yaml:"ssl_client_key"`
	JwtSecret     string `yaml:"jwt_secret"`

	HandshakePort int `yaml:"handshake_port"`
	CommonPort    int `yaml:"common_port"`
	BroadcastPort int `yaml:"broadcast_port"`

	DatabaseFile string `yaml:"database_file"`
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func makeConfigRelativePath(configPath string, relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	return filepath.Join(filepath.Dir(configPath), relativePath)
}

func LoadConfig(path string) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.JwtSecret = makeConfigRelativePath(path, cfg.JwtSecret)
	cfg.SslCaCert = makeConfigRelativePath(path, cfg.SslCaCert)
	cfg.SslClientCert = makeConfigRelativePath(path, cfg.SslClientCert)
	cfg.SslClientKey = makeConfigRelativePath(path, cfg.SslClientKey)
	cfg.DatabaseFile = makeConfigRelativePath(path, cfg.DatabaseFile)

	if !fileExists(cfg.JwtSecret) {
		return nil, fmt.Errorf("jwt secret file does not exist: %s", cfg.JwtSecret)
	}
	if !fileExists(cfg.DatabaseFile) {
		return nil, fmt.Errorf("database file does not exist: %s", cfg.DatabaseFile)
	}
	if !fileExists(cfg.SslCaCert) {
		return nil, fmt.Errorf("ssl ca cert file does not exist: %s", cfg.SslCaCert)
	}
	if !fileExists(cfg.SslClientCert) {
		return nil, fmt.Errorf("ssl client cert file does not exist: %s", cfg.SslClientCert)
	}
	if !fileExists(cfg.SslClientKey) {
		return nil, fmt.Errorf("ssl client key file does not exist: %s", cfg.SslClientKey)
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
