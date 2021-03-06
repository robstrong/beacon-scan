package main

import (
	"flag"
	"log"
	"math"

	"github.com/eclipse/paho.mqtt.golang"
)

var configPath = flag.String("config", "config.toml", "path to config file")

func main() {
	flag.Parse()

	//load config
	c := loadConfig()

	//setup BLE scanner and MQTT client
	s := NewBLEScanner()
	mq := getMQTTClient(c.MQTT)

	//register each UUID from config
	for _, u := range c.UUIDs {
		log.Printf("adding %s", u.UUID)
		s.AddSubscriber(handleBeaconAdvert(u, mq))
	}

	//start listening
	if err := s.Start(); err != nil {
		log.Fatalf("error starting scanner: %v", err)
	}
}

func calculateDistance(rssi int, txPower int) float64 {
	if rssi == 0 || txPower == 0 {
		return -1
	}

	var ratio = float64(rssi) / float64(txPower)
	if ratio < 1 {
		return math.Pow(ratio, 10)
	}
	return (0.89976)*math.Pow(ratio, 7.7095) + 0.111
}

func handleBeaconAdvert(u TrackedUUID, mq mqtt.Client) BeaconDiscovered {
	return func(d BeaconData) {
		if d.UUID != u.UUID {
			return
		}
		payload := newPayload(u.DeviceID, u.Name, calculateDistance(d.RSSI, d.TxPower))
		log.Printf("publishing to %s: %s\n", u.Topic, payload)
		t := mq.Publish(u.Topic, 0, false, payload)
		t.Wait()
		if err := t.Error(); err != nil {
			log.Printf("error publishing: %v", err)
		}
	}
}
