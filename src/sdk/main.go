package sdk

/*
 * TODO:
 * 1. Possible race-condition when connecting to peripheral (very small probability)
 * 2. Error handling in http (prot package)
 * 
 */

import (
	"clients/prot"
	"db"
	"typeutil"
	"ble"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type MoecoSDK struct {
	db                      *db.DBAdapter
	client                  *prot.Client
	host                    string
	apiKey                  string
	gatewayHash             string
	dbPath                  string
	lastSync                int
	getDevicesInterval      int
	syncInterval            int
	charNotifyInterval      int
	deviceConnInterval      int
	transactionsBufSize     int
	stoped                  bool
	errors                  *chan error
	transactions            chan db.Transaction
	deviceTimeouts          map[string]time.Time
	ble                     *ble.MoecoBLE
	log                     *logrus.Logger
}

func NewMoecoSDK(host, apiKey, gatewayHash, dbPath string) MoecoSDK {
	return MoecoSDK{
		host:                    host,
		apiKey:                  apiKey,
		gatewayHash:             gatewayHash,
		dbPath:                  dbPath,
		getDevicesInterval:      10000000,
		syncInterval:            4000000,
		charNotifyInterval:      5000000,
		deviceConnInterval:      60000000,
		transactionsBufSize:     50,
	}
}

func (m *MoecoSDK) Start(log *logrus.Logger) (error, chan error) {
	errorsChan := make(chan error)
	sqliteDb, err := db.NewDBAdapter(m.dbPath)
	if err != nil {
		return errors.Wrap(err, "db adapter init failed"), nil
	}
	client := prot.NewClient(m.host, m.apiKey, m.gatewayHash)
	err = client.Init(log)
	if err != nil {
		return errors.Wrap(err, "gateway client init failed"), nil
	}

	m.transactions = make(chan db.Transaction, m.transactionsBufSize)
	ble, err := ble.NewMoecoBLE(log, sqliteDb, &errorsChan, &m.transactions,
	m.transactionsBufSize, m.charNotifyInterval, m.deviceConnInterval)
	if err != nil {
		return errors.Wrap(err, "MoecoBLE init failed"), nil
	}

	m.db = sqliteDb
	m.client = &client
	m.errors = &errorsChan
	m.ble = ble
	m.log = log
	go m.getTransactions()
	go m.runSync()
	go m.getDevices()
	return nil, *m.errors
}

func (m *MoecoSDK) getTransactions() {
	for {
		t := <-m.transactions
		m.log.Debugf("Add transaction: %s", t)
		err := m.db.InsertTransaction(t)
		if err != nil {
			*m.errors <- errors.Wrap(err, "insert transaction failed")
		}
	}
}

func (m *MoecoSDK) runSync() {
	for range time.Tick(time.Duration(m.syncInterval) * time.Microsecond) {
		if m.stoped {
			break
		}
		m.log.Info("Sync transactions")
		temp, err := m.db.GetUnsendTransaction()
		if err != nil {
			*m.errors <- errors.Wrap(err, "getting unsend transactions failed")
			continue
		}
		if len(temp) == 0 {
			continue
		}
		tr := types.TrasactionsToReq(temp)
		lastSync := 0
		ids := make([]int, 0, len(tr))
		for _, v := range tr {
			if int(v.Timestamp.Unix()) > lastSync {
				lastSync = int(v.Timestamp.Unix())
			}
			ids = append(ids, v.ID)
		}
		_, err = m.client.SyncTransaction(prot.Transactions{
			Transactions: tr,
		}, m.lastSync)
		if err != nil {
			*m.errors <- errors.Wrap(err, "transactions sync failed")
			continue
		}
		err = m.db.SetSendedTransaction(ids)
		if err != nil {
			*m.errors <- errors.Wrap(err, "set send status on transactions failed")
			continue
		}
		m.lastSync = int(time.Now().Unix())
	}
}

func (m *MoecoSDK) getDevices() {
	for range time.Tick(time.Duration(m.getDevicesInterval) * time.Microsecond) {
		if m.stoped {
			break
		}
		m.log.Info("Get devices")
		res, err := m.client.GetDevices()
		if err != nil {
			*m.errors <- errors.Wrap(err, "get devices failed")
			continue
		}
		if len(res.Data) == 0 {
			*m.errors <- fmt.Errorf("invalid response get device")
			continue
		}
		deviceGroups, err := types.DeviceGroupsFromResponse(res.Data[0].DeviceGroups)
		if err != nil {
			*m.errors <- err
			continue
		}
		err = m.db.InsertDeviceGroups(deviceGroups)
		if err != nil {
			*m.errors <- errors.Wrap(err, "device groups db insertion failed")
			continue
		}
		devices := types.DevicesFromResponse(res.Data[0].Devices)
		err = m.db.InsertDevices(devices)
		if err != nil {
			*m.errors <- errors.Wrap(err, "devices db insertion failed")
			continue
		}
	}
}
