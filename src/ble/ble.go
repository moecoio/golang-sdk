package ble

import (
	"clients/prot"
	"typeutil"
	"db"
	"log"
	"fmt"
	"strings"
	"encoding/json"
	"encoding/hex"
	"time"

	"github.com/mihalicyn/gatt"
	"github.com/mihalicyn/gatt/examples/option"
	"github.com/sirupsen/logrus"
)

type MoecoBLE struct {
	log                     *logrus.Logger
	db                      *db.DBAdapter
	charNotifyInterval      int
	deviceConnInterval      int
	transactionsBufSize     int
	errors                  *chan error
	transactions            *chan db.Transaction
	deviceTimeouts          map[string]time.Time
}

func NewMoecoBLE(
	logger *logrus.Logger, database *db.DBAdapter,errors *chan error,
	transactions *chan db.Transaction, transactionsBufSize int,
	charNotifyInterval int, deviceConnInterval int) (*MoecoBLE, error) {
	// catch logs from gatt
	log.SetOutput(logger.Writer())

	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		logger.Fatalf("Failed to open device, err: %s\n", err)
		return nil, err
	}

	//m.transactions = make(chan db.Transaction, m.transactionsBufSize)
	deviceTimeouts := make(map[string]time.Time)

	ble := &MoecoBLE{
		log:                 logger,
		db:                  database,
		charNotifyInterval:  charNotifyInterval,
		deviceConnInterval:  deviceConnInterval,
		transactionsBufSize: transactionsBufSize,
		errors:              errors,
		transactions:        transactions,
		deviceTimeouts:      deviceTimeouts,
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(genOnPeriphDiscoveredCbk(ble)),
		gatt.PeripheralConnected(genOnPeriphConnectedCbk(ble)),
		gatt.PeripheralDisconnected(genOnPeriphDisconnectedCbk(ble)),
	)

	d.Init(func(d gatt.Device, s gatt.State) {
		switch s {
		case gatt.StatePoweredOn:
			d.Scan([]gatt.UUID{}, true)
			return
		default:
			d.StopScanning()
		}
	})

	return ble, nil
}

func genOnPeriphDiscoveredCbk(ble *MoecoBLE) func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	return func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
		ble.log.Debugf("\nFound... Peripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())

		devices, err := ble.db.GetDevices()
		if err != nil {
			*ble.errors <- err
			return
		}

		for _, device := range devices {
			if strings.EqualFold(p.ID(), device.Hash) {
				// not connect to device before timeout
				timeout, ok := ble.deviceTimeouts[p.ID()]
				if ok && timeout.After(time.Now()) {
					return
				}

				ble.log.Infof("Peripheral found in whitelist ID: %s, Name: %s\n", p.ID(), p.Name())
				// Stop scanning once we've got one of the peripherals we're looking for.
				p.Device().StopScanning()

				// now we can connect
				p.Device().Connect(p)
				return
			}
		}

		ble.log.Debugf("Peripheral not found in whitelist ID: %s, Name: %s\n", p.ID(), p.Name())
	}
}

func genOnPeriphConnectedCbk(ble *MoecoBLE) func(p gatt.Peripheral, err error) {
	return func(p gatt.Peripheral, err error) {
		ble.log.Infof("Connected to %s %s\n", p.ID(), p.Name())
		defer p.Device().CancelConnection(p)

		if err := p.SetMTU(500); err != nil {
			*ble.errors <- fmt.Errorf("failed to set MTU, err: %s\n", err)
		}

		device, err := ble.db.GetDeviceByHash(p.ID())
		if err != nil {
			*ble.errors <- err
			return
		}
		if device == nil {
			*ble.errors <- fmt.Errorf("already connected device not found in whitelist ID: %s\n", p.ID())
			return
		}

		deviceGroupDB, err := ble.db.GetDeviceGroupByID(device.DeviceGroupID)
		if err != nil {
			*ble.errors <- err
			return
		}
		if deviceGroupDB == nil {
			*ble.errors <- fmt.Errorf("device group not found for device with ID: %s, device group ID: %s\n", p.ID(), device.DeviceGroupID)
			return
		}

		deviceGroup, err := types.DeviceGroupToResponse(*deviceGroupDB)
		if err != nil {
			*ble.errors <- err
			return
		}

		// Discovery device services
		pServices, err := p.DiscoverServices(nil)
		if err != nil {
			ble.log.Errorf("failed to discover services, err: %s\n", err)
			return
		}

		payload := make(map[string]map[string]string)

		for _, pService := range pServices {
			var dgService *prot.Service = nil
			for _, s := range deviceGroup.Services {
				name := strings.Replace(s.Name, "-", "", -1)
				if strings.EqualFold(pService.UUID().String(), name) {
					dgService = &s
					break
				}
			}
			if dgService == nil {
				continue
			}

			// Discovery characteristics
			cs, err := p.DiscoverCharacteristics(nil, pService)
			if err != nil {
				ble.log.Errorf("failed to discover characteristics, err: %s\n", err)
				continue
			}

			for _, pChar := range cs {
				var dgChar *prot.Characteristic
				for _, c := range dgService.Characteristics {
					name := strings.Replace(c.Name, "-", "", -1)
					if strings.EqualFold(pChar.UUID().String(), name) {
						dgChar = &c
						break
					}
				}
				if dgChar == nil {
					continue
				}

				// Read the characteristic, if possible.
				if (pChar.Properties() & gatt.CharRead) != 0 {
					b, err := p.ReadLongCharacteristic(pChar)
					if err != nil {
						ble.log.Errorf("failed to read characteristic, err: %s\n", err)
						continue
					}

					// init map if is not yet allocated
					_, ok := payload[dgService.Name]
					if !ok {
						payload[dgService.Name] = make(map[string]string)
					}
					payload[dgService.Name][dgChar.Name] = hex.EncodeToString(b)
				}

				/*
				 * It's needed to discover descriptors *before*
				 * attempting to subscribe to characteristic
				 */
				_, err := p.DiscoverDescriptors(nil, pChar)
				if err != nil {
					ble.log.Warn("failed to discover descriptors, err: %s\n", err)
					continue
				}

				// Subscribe the characteristic, if possible.
				if (pChar.Properties() & (gatt.CharNotify | gatt.CharIndicate)) != 0 {
					f := func(c *gatt.Characteristic, b []byte, err error) {
						// init map if is not yet allocated
						_, ok := payload[dgService.Name]
						if !ok {
							payload[dgService.Name] = make(map[string]string)
						}
						payload[dgService.Name][dgChar.Name] = hex.EncodeToString(b)
					}
					if err := p.SetNotifyValue(pChar, f); err != nil {
						ble.log.Errorf("failed to subscribe characteristic, err: %s\n", err)
						continue
					}
				}
			}
		}

		// Waiting to get some notifiations, if any.
		time.Sleep(time.Duration(ble.charNotifyInterval) * time.Microsecond)

		// Prepare transaction and put it on channel
		b, err := json.Marshal(payload)
		ble.log.Debugf("payload: %s\n", b);
		*ble.transactions <- db.Transaction{
			Hash:       "",
			DeviceHash: device.Hash,
			Timestamp:  int(time.Now().Unix()), // FIXME: truncating int64->int32
			Uplink:     0,
			Sended:     0,
			Payload:    string(b),
		}
	}
}

func genOnPeriphDisconnectedCbk(ble *MoecoBLE) func(p gatt.Peripheral, err error) {
	return func(p gatt.Peripheral, err error) {
		// init device timeout
		ble.deviceTimeouts[p.ID()] = time.Now().Add(time.Duration(ble.deviceConnInterval) * time.Microsecond)

		//delete(m.deviceTimeouts, p.ID())
		ble.log.Infof("Disconnected from %s\n", p.ID())
		// turn on scan
		p.Device().Scan([]gatt.UUID{}, true)
	}
}

