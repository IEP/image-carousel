package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holder that consists of bot token based on the platform
type Config struct {
	TelegramToken string `mapstructure:"telegram_token"`
	DiscordToken  string `mapstructure:"discord_token"`
}

// LoadConfig from provided config name
func LoadConfig(name string) (*Config, error) {
	viper.SetConfigName(name)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper.ReadInConfig: %w", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("viper.Unmarshal: %w", err)
	}
	return &cfg, nil
}
