package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Telegram struct {
		Token  string `yaml:"token"`
		ChatID string `yaml:"chat_id"`
	} `yaml:"telegram"`
	Tracking struct {
		URL           string   `yaml:"url"`
		Interval      string   `yaml:"interval"`
		SteaksToTrack []string `yaml:"steaks_to_track"`
	} `yaml:"tracking"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
