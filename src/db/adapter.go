package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

const (
	createDeviceTable = "CREATE TABLE IF NOT EXISTS device(" +
		"id              INTEGER PRIMARY KEY," +
		"hash            TEXT UNIQUE," +
		"manufacturer    TEXT," +
		"created_at      INTEGER," +
		"updated_at      INTEGER," +
		"exonum_id       TEXT," +
		"device_group_id TEXT," +
		"owner_key       TEXT" +
		")"
	createDeviceGroupTable = "CREATE TABLE IF NOT EXISTS device_group(" +
		"id                INTEGER PRIMARY KEY," +
		"exonum_id         TEXT UNIQUE," +
		"name              TEXT," +
		"group_type        INTEGER," +
		"uplink_lifetime   INTEGER," +
		"downlink_lifetime INTEGER," +
		"services          TEXT," +
		"created_at        INTEGER," +
		"updated_at        INTEGER," +
		"owner_key         TEXT" +
		")"
	createTransactionTable = "CREATE TABLE IF NOT EXISTS tr(" +
		"id INTEGER PRIMARY KEY, " +
		"hash TEXT, " +
		"device_hash TEXT," +
		"timestamp INTEGER," +
		"uplink INTEGER," +
		"sended INTEGER," +
		"payload TEXT" +
		")"

	transactionInsertQuery = "INSERT INTO tr " +
		"(hash, device_hash, timestamp, uplink, sended, payload) " +
		"VALUES (?, ?, ?, ?, ?, ?)"
	deviceInsertQuery = "INSERT INTO device " +
		"(hash, manufacturer, created_at, updated_at, exonum_id, device_group_id, owner_key) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)" +
		"ON CONFLICT(hash) DO " +
		"UPDATE SET " +
		"manufacturer = $2, created_at = $3, updated_at = $4, exonum_id = $5, " +
		"device_group_id = $6, owner_key = $7 " +
		"WHERE hash = $1"
	deviceGroupInsertQuery = "INSERT INTO device_group " +
		"(" +
		"exonum_id, name, group_type, uplink_lifetime, downlink_lifetime, " +
		"services, created_at, updated_at,  owner_key " +
		") " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)" +
		"ON CONFLICT(exonum_id) DO " +
		"UPDATE SET " +
		"name = $2, group_type = $3, uplink_lifetime = $4, downlink_lifetime = $5, " +
		"services = $6, created_at = $7, updated_at = $8, owner_key = $9 " +
		"WHERE exonum_id = $1"
	deviceGetQuery = "SELECT " +
		"id, hash, manufacturer, created_at, updated_at, exonum_id, device_group_id, owner_key " +
		"FROM device"
	deviceGetByHashQuery = "SELECT " +
		"id, hash, manufacturer, created_at, updated_at, exonum_id, device_group_id, owner_key " +
		"FROM device WHERE LOWER(hash) = LOWER($1)"
	deviceGroupGetQuery = "SELECT " +
		"id, exonum_id, name, group_type, uplink_lifetime, downlink_lifetime, " +
		"services, created_at, updated_at,  owner_key " +
		"FROM device_group"
	deviceGroupGetByIdQuery = "SELECT " +
		"id, exonum_id, name, group_type, uplink_lifetime, downlink_lifetime, " +
		"services, created_at, updated_at,  owner_key " +
		"FROM device_group WHERE LOWER(exonum_id) = LOWER($1)"
	transactionGetQuery = "SELECT " +
		"id, hash, device_hash, timestamp, uplink, sended, payload " +
		"FROM tr WHERE sended = 0"
)



type DBAdapter struct {
	db                    *sql.DB
	transactionInsertStmt *sql.Stmt
	deviceGroupInsertStmt *sql.Stmt
	deviceInsertStmt      *sql.Stmt
}

func NewDBAdapter(path string) (*DBAdapter, error) {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = database.Exec(createDeviceTable)
	if err != nil {
		return nil, err
	}
	_, err = database.Exec(createDeviceGroupTable)
	if err != nil {
		return nil, err
	}
	_, err = database.Exec(createTransactionTable)
	if err != nil {
		return nil, err
	}
	transactionInsertStmt, err := database.Prepare(transactionInsertQuery)
	if err != nil {
		return nil, err
	}
	deviceInsertStmt, err := database.Prepare(deviceInsertQuery)
	if err != nil {
		return nil, err
	}
	deviceGroupInsertStmt, err := database.Prepare(deviceGroupInsertQuery)
	if err != nil {
		return nil, err
	}
	return &DBAdapter{
		db:                    database,
		transactionInsertStmt: transactionInsertStmt,
		deviceInsertStmt:      deviceInsertStmt,
		deviceGroupInsertStmt: deviceGroupInsertStmt,
	}, nil
}

func (db *DBAdapter) InsertTransaction(v Transaction) error {
	_, err := db.transactionInsertStmt.Exec(
		v.Hash,
		v.DeviceHash,
		v.Timestamp,
		v.Uplink,
		v.Sended,
		v.Payload)

	return err
}

func (db *DBAdapter) InsertTransactions(res []Transaction) error {
	for _, v := range res {
		_, err := db.transactionInsertStmt.Exec(
			v.Hash,
			v.DeviceHash,
			v.Timestamp,
			v.Uplink,
			v.Sended,
			v.Payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DBAdapter) InsertDevices(devices []Device) error {
	for _, device := range devices {
		_, err := db.deviceInsertStmt.Exec(
			device.Hash,
			device.Manufacturer,
			device.CreatedAt,
			device.UpdatedAt,
			device.ExonumID,
			device.DeviceGroupID,
			device.OwnerKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DBAdapter) InsertDeviceGroups(deviceGroups []DeviceGroup) error {
	for _, deviceGroup := range deviceGroups {
		_, err := db.deviceGroupInsertStmt.Exec(
			deviceGroup.ExonumID,
			deviceGroup.Name,
			deviceGroup.GroupType,
			deviceGroup.UplinkLifetime,
			deviceGroup.DownlinkLifetime,
			deviceGroup.Services,
			deviceGroup.CreatedAt,
			deviceGroup.UpdatedAt,
			deviceGroup.OwnerKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DBAdapter) SetSendedTransaction(ids []int) error {
	where := make([]string, 0, len(ids))
	for _, v := range ids {
		where = append(where, "id = "+strconv.Itoa(v))
	}
	_, err := db.db.Exec("UPDATE tr " +
		"SET sended = 1 " +
		"WHERE " + strings.Join(where, " OR "))
	if err != nil {
		return err
	}
	return nil
}

func (db *DBAdapter) GetDevices() ([]Device, error) {
	rows, err := db.db.Query(deviceGetQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var devices []Device
	for rows.Next() {
		var d Device
		err = rows.Scan(
			&d.ID,
			&d.Hash,
			&d.Manufacturer,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.ExonumID,
			&d.DeviceGroupID,
			&d.OwnerKey)
		if err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func (db *DBAdapter) GetDeviceByHash(hash string) (*Device, error) {
	rows, err := db.db.Query(deviceGetByHashQuery, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d Device
		err = rows.Scan(
			&d.ID,
			&d.Hash,
			&d.Manufacturer,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.ExonumID,
			&d.DeviceGroupID,
			&d.OwnerKey)
		if err != nil {
			return nil, err
		}
		return &d, nil
	}
	return nil, nil
}

func (db *DBAdapter) GetDeviceGroups() ([]DeviceGroup, error) {
	rows, err := db.db.Query(deviceGroupGetQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var deviceGroups []DeviceGroup
	for rows.Next() {
		var d DeviceGroup
		err = rows.Scan(
			&d.ID,
			&d.ExonumID,
			&d.Name,
			&d.GroupType,
			&d.UplinkLifetime,
			&d.DownlinkLifetime,
			&d.Services,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.OwnerKey)
		if err != nil {
			return nil, err
		}
		deviceGroups = append(deviceGroups, d)
	}
	return deviceGroups, nil
}

func (db *DBAdapter) GetDeviceGroupByID(id string) (*DeviceGroup, error) {
	rows, err := db.db.Query(deviceGroupGetByIdQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d DeviceGroup
		err = rows.Scan(
			&d.ID,
			&d.ExonumID,
			&d.Name,
			&d.GroupType,
			&d.UplinkLifetime,
			&d.DownlinkLifetime,
			&d.Services,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.OwnerKey)
		if err != nil {
			return nil, err
		}
		return &d, nil
	}
	return nil, nil
}

func (db *DBAdapter) GetUnsendTransaction() ([]Transaction, error) {
	rows, err := db.db.Query(transactionGetQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		rows.Scan(&t.ID, &t.Hash, &t.DeviceHash, &t.Timestamp, &t.Uplink, &t.Sended, &t.Payload)
		transactions = append(transactions, t)
	}
	return transactions, nil

}
