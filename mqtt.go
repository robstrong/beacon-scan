package main

import (
	"encoding/json"
	"log"
	"time"

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
	log.Printf("connecting to MQTT at %s", config.Host)
	opt := mqtt.NewClientOptions().
		AddBroker(config.Host).
		SetAutoReconnect(true).
		SetCredentialsProvider(mqtt.CredentialsProvider(func() (string, string) {
			return config.Username, config.Password
		})).
		SetKeepAlive(2 * time.Second).
		SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opt)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c
}
