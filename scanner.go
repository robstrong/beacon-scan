package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type BeaconData struct {
	UUID    string
	Major   uint16
	Minor   uint16
	RSSI    int
	TxPower int
}

func parseBeaconData(data []byte) (*BeaconData, error) {
	if len(data) < 25 || binary.BigEndian.Uint32(data) != 0x4c000215 {
		return nil, errors.New("Not an iBeacon")
	}
	beacon := new(BeaconData)
	beacon.UUID = strings.ToUpper(hex.EncodeToString(data[4:8]) + "-" + hex.EncodeToString(data[8:10]) + "-" + hex.EncodeToString(data[10:12]) + "-" + hex.EncodeToString(data[12:14]) + "-" + hex.EncodeToString(data[14:20]))
	beacon.Major = binary.BigEndian.Uint16(data[20:22])
	beacon.Minor = binary.BigEndian.Uint16(data[22:24])
	return beacon, nil
}

func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		device.Scan(nil, true)
		return
	default:
		device.StopScanning()
	}
}

type BeaconDiscovered func(BeaconData)

func (b *BLEScanner) onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	data, err := parseBeaconData(a.ManufacturerData)
	if err != nil {
		return
	}
	data.RSSI = rssi
	data.TxPower = -59
	if a.TxPowerLevel != 0 {
		data.TxPower = a.TxPowerLevel
	}
	b.subscribersMu.RLock()
	defer b.subscribersMu.RUnlock()
	for _, s := range b.subscribers {
		s(*data)
	}
}

type BLEScanner struct {
	device        gatt.Device
	subscribers   []BeaconDiscovered
	subscribersMu sync.RWMutex
}

func NewBLEScanner() *BLEScanner {
	return &BLEScanner{}
}

func (b *BLEScanner) Start() error {
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		return fmt.Errorf("Failed to open device, err: %s\n", err)
	}
	device.Handle(gatt.PeripheralDiscovered(b.onPeripheralDiscovered))
	if err := device.Init(onStateChanged); err != nil {
		return err
	}
	b.device = device
	select {}
}

func (b *BLEScanner) AddSubscriber(fn BeaconDiscovered) {
	b.subscribersMu.Lock()
	b.subscribers = append(b.subscribers, fn)
	b.subscribersMu.Unlock()
}
