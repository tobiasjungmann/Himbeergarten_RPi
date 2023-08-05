package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

func handleBluetooth() {
	log.Fatalf("Bluetooth forwarding not yet implemented")
}
func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Printf("State: %s\n", s)
	switch s {
	case gatt.StatePoweredOn:
		// Start scanning for nearby BLE devices
		fmt.Println("Scanning...")
		d.Scan([]gatt.UUID{}, true)
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	// Print the name of the discovered peripheral
	if len(p.Name()) > 0 {
		fmt.Printf("Peripheral discovered: %s\n", p.Name())
	}
}

func BluetoothTest() {
	device, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		fmt.Println("Failed to open the default Bluetooth adapter:", err)
		return
	}

	// Register event handlers
	device.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
	)

	// Start the device
	device.Init(onStateChanged)
	defer device.StopScanning()

	// Wait for a few seconds to scan for devices
	time.Sleep(5 * time.Second)
}

// todo arduino starten, der werte über bluetooth sendet - diese hier empfangen und erstaml nur ausdrucken
/*
1. Handshake durchführen, damit die verbindung überhaupt zustande kommt
2. Daten in ein format bringen, in dem sie weiterverarbeitet werden können
*/
/*
func BluetoothTest() {
	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s", err)
	}

	// Register optional handlers.
	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) { fmt.Println("Connect: ", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) { fmt.Println("Disconnect: ", c.ID()) }),
	)

	// A mandatory handler for monitoring device state.
	onStateChanged := func(d gatt.Device, s gatt.State) {
		fmt.Printf("State: %s\n", s)
		switch s {
		case gatt.StatePoweredOn:
			// Setup GAP and GATT services for Linux implementation.
			// OS X doesn't export the access of these services.
			d.AddService(service.NewGapService("Gopher")) // no effect on OS X
			d.AddService(service.NewGattService())        // no effect on OS X

			// A simple count service for demo.
			s1 := service.NewCountService()
			d.AddService(s1)

			// A fake battery service for demo.
			s2 := service.NewBatteryService()
			d.AddService(s2)

			// Advertise device name and service's UUIDs.
			d.AdvertiseNameAndServices("Gopher", []gatt.UUID{s1.UUID(), s2.UUID()})

			// Advertise as an OpenBeacon iBeacon
			d.AdvertiseIBeacon(gatt.MustParseUUID("AA6062F098CA42118EC4193EB73CCEB6"), 1, 2, -59)

		default:
		}
	}

	d.Init(onStateChanged)
	select {}
}
*/
