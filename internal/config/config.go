package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// LoggerConfig stores logger level and file if needed to log into
type LoggerConfig struct {
	Level string `default:"debug" env:"LOG_LEVEL" yaml:"level"`
	File  string `default:""      env:"LOG_FILE"  yaml:"file"`
}

// Cert ssl certificate config
type Cert struct {
	Cert string `yaml:"cert" env:"CERT_FILE"`
	Key  string `yaml:"key" env:"CERT_KEY"`
}

type S3Config struct {
	Endpoint     string `yaml:"endpoint" env:"DAD_S3_ENDPOINT"`
	Region       string `yaml:"region" env:"DAD_S3_REGION"`
	ExpiresURLIn int    `yaml:"expires_url_in" env:"S3_EXPIRES_URL"`
	AccessKey    string `yaml:"access_key" env:"DAD_S3_ACCESS_KEY"`
	SecretKey    string `yaml:"secret_key" env:"DAD_S3_SECRET_KEY"`
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
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// DadConfig stores all config for the application
type DadConfig struct {
	LoggerCfg     LoggerConfig     `yaml:"logger"`
	GRPSServerCfg GRPSServerConfig `yaml:"grps_server"`
	Cert          Cert             `yaml:"cert"`
	DBCfg         DBConfig         `yaml:"db"`
	S3Cfg         S3Config         `yaml:"s3"`
}

// GetGRPSAddress returns address:port of grps server
func (d *DadConfig) GetGRPSAddress() string {
	return fmt.Sprintf("%s:%d", d.GRPSServerCfg.Host, d.GRPSServerCfg.Port)
}

type GRPSServerConfig struct {
	Host                 string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
	Port                 int    `yaml:"port" env:"SERVER_PORT" env-default:"8081"`
	Secure               bool   `yaml:"secure" env:"SERVER_SECURE" env-default:"false"`
	AccessSecret         string `yaml:"access_secret" env:"SERVER_ACCESS_SECRET" env-default:""`
	RefreshSecret        string `yaml:"refresh_secret" env:"SERVER_REFRESH_SECRET" env-default:""`
	AccessTokenDuration  int    `yaml:"access_token_duration" env:"ACCESS_TOKEN_DURATION" env-default:"15"`
	RefreshTokenDuration int    `yaml:"refresh_token_duration" env:"REFRESH_TOKEN_DURATION" env-default:"43200"`
}

// NewDadConfig is a constructor for DadConfig
func NewDadConfig(file string) (*DadConfig, error) {

	_ = godotenv.Load()

	config := &DadConfig{}
	if err := cleanenv.ReadConfig(file, config); err != nil {
		return nil, fmt.Errorf("unable to read application config: %w", err)
	}
	return config, nil
}
