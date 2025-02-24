package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

// LoggerConfig stores logger level and file if needed to log into
type LoggerConfig struct {
	Level string `default:"debug" env:"LOG_LEVEL" yaml:"level"`
	File  string `default:""      env:"LOG_FILE"  yaml:"file"`
}

type Cert struct {
	Cert string `yaml:"cert" env:"CERT_FILE"`
	Key  string `yaml:"key" env:"CERT_KEY"`
}

// DadConfig stores all config for the application
type DadConfig struct {
	LoggerCfg     LoggerConfig     `yaml:"logger"`
	GRPSServerCfg GRPSServerConfig `yaml:"grps_server"`
	Cert          Cert             `yaml:"cert"`
}

// GetGRPSAddress returns address:port of grps server
func (d *DadConfig) GetGRPSAddress() string {
	return fmt.Sprintf("%s:%d", d.GRPSServerCfg.Host, d.GRPSServerCfg.Port)
}

type GRPSServerConfig struct {
	Host          string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
	Port          int    `yaml:"port" env:"SERVER_PORT" env-default:"8081"`
	Secure        bool   `yaml:"secure" env:"SERVER_SECURE" env-default:"false"`
	Secret        string `yaml:"secret" env:"SERVER_SECRET" env-default:""`
	TokenDuration int    `yaml:"token_duration" env:"TOKEN_DURATION" env-default:"300"`
}

// NewDadConfig is a constructor for DadConfig
func NewDadConfig(file string) (*DadConfig, error) {
	config := &DadConfig{}

	if err := cleanenv.ReadConfig(file, config); err != nil {
		return nil, fmt.Errorf("unable to read application config: %w", err)
	}
	return config, nil
}
