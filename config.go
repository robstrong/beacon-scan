package main

import (
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

func loadConfig() Config {
	c := Config{}
	_, err := toml.DecodeFile(*configPath, &c)
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}
	return normalizeConfig(c)
}

type Config struct {
	UUIDs []TrackedUUID
	MQTT  MQTTConfig
}

type TrackedUUID struct {
	Name     string
	UUID     string
	Topic    string
	DeviceID string
}

type MQTTConfig struct {
	Host     string
	Username string
	Password string
}

func normalizeConfig(c Config) Config {
	for i, t := range c.UUIDs {
		c.UUIDs[i].UUID = strings.ToUpper(t.UUID)
	}
	return c
}
