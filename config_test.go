package main

import "testing"

func TestConfig(t *testing.T) {
	cfg, err := LoadConfig("config.example")
	if err != nil {
		t.Errorf("LoadConfig: %s", err)
	}

	if cfg.DiscordToken != "dt" || cfg.TelegramToken != "tt" {
		t.Errorf("invalid config; config: %+v", cfg)
	}

	t.Logf("config: %+v", cfg)
}
