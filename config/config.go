package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DSN      string `mapstructure:"dsn"`
	LogLevel string `mapstructure:"log_level"`
	Server   Server `mapstructure:"server"`
	Kafka    Kafka  `mapstructure:"kafka"`
}

type Server struct {
	Address string `mapstructure:"address"`
}

type Kafka struct {
	Brokers string `mapstructure:"brokers"`
	Topic   string `mapstructure:"topic"`
	GroupID string `mapstructure:"group_id"`
}

var Cfg *Config

func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
