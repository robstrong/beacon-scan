package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

func loadConfig() Config {
	c := Config{}
	_, err := toml.DecodeFile("config.toml", &c)
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}
	return c
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
