package main

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type payload struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Distance float64 `json:"distance"`
}

func newPayload(id, name string, distance float64) []byte {
	return mustMarshalJSON(&payload{
		ID:       id,
		Name:     name,
		Distance: distance,
	})
}

func mustMarshalJSON(p *payload) []byte {
	v, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return v
}
func getMQTTClient(config MQTTConfig) mqtt.Client {
	opt := mqtt.NewClientOptions().
		AddBroker(config.Host).
		SetAutoReconnect(true).
		SetCredentialsProvider(mqtt.CredentialsProvider(func() (string, string) {
			return config.Username, config.Password
		}))
	c := mqtt.NewClient(opt)
	t := c.Connect()
	if err := t.Error(); err != nil {
		log.Fatalf("error connecting to mqtt broker: %v", err)
	}
	return c
}
