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

// DBConfig stores database creds and addresses
type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_AME" env-default:"owl"`
	User     string `yaml:"username"  env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSL_MODE" env-default:"disable"`
}

func (d *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.User, d.Password, d.Name, d.SSLMode)
}

// DadConfig stores all config for the application
type DadConfig struct {
	LoggerCfg     LoggerConfig     `yaml:"logger"`
	GRPSServerCfg GRPSServerConfig `yaml:"grps_server"`
	Cert          Cert             `yaml:"cert"`
	DBCfg         DBConfig         `yaml:"db"`
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
