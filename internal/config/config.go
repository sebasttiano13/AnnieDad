package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type LoggerConfig struct {
	Level string `default:"debug" env:"LOG_LEVEL" yaml:"level"`
	File  string `default:""      env:"LOG_FILE"  yaml:"file"`
}

type ServerConfig struct {
	LoggerCfg LoggerConfig `yaml:"logger"`
}
