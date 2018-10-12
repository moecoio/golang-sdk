package prot

import (
	"encoding/json"
	"time"
)

type Meta struct {
	Total  int         `json:"total"`
	Count  int         `json:"count"`
	Offset int         `json:"offset"`
	Error  interface{} `json:"error"`
}

type BaseResponse struct {
	Meta Meta            `json:"meta"`
	Data json.RawMessage `json:"data"`
}

type SyncResponse struct {
	BaseResponse
	Data []SyncResponseData `json:"data"`
}

type GateInitResponse struct {
	BaseResponse
	Data []GateInitResponseData `json:"data"`
}

type GateInitResponseData struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Hash      string    `json:"hash"`
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
}

type TransactionRes struct {
	ID              int        `json:"id"`
	DeviceHash      string     `json:"device_hash"`
	GatewayID       int        `json:"gateway_id"`
	Hash            string     `json:"hash"`
	Timestamp       time.Time  `json:"timestamp"`
	Status          int        `json:"status"`
	Uplink          bool       `json:"uplink"`
	Payload         string     `json:"payload"`
	RejectionReason *string    `json:"rejection_reason"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ExpireDate      *time.Time `json:"expire_date"`
	InvoiceID       *string    `json:"invoice_id"`
}

type SyncResponseData struct {
	UpdatedUplink []TransactionRes `json:"updatedUplink"`
	Changed       []TransactionRes `json:"changed"`
	Uplink        []TransactionRes `json:"uplink"`
	Results       []TransactionRes `json:"results"`
}

type DeviceGroup struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	GroupType        int       `json:"group_type"`
	UplinkLifetime   int       `json:"uplink_lifetime"`
	DownlinkLifetime int       `json:"downlink_lifetime"`
	Services         []Service `json:"services"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ExonumID         string    `json:"exonum_id"`
	OwnerKey         string    `json:"owner_key"`
}

type Service struct {
	Name            string           `json:"name"`
	Characteristics []Characteristic `json:"characteristics"`
}

type Characteristic struct {
	Mac        bool   `json:"mac"`
	Name       string `json:"name"`
	Readable   bool   `json:"readable"`
	Writable   bool   `json:"writable"`
	Notifiable bool   `json:"notifiable"`
}

type Device struct {
	ID            string    `json:"id"`
	Hash          string    `json:"hash"`
	Manufacturer  string    `json:"manufacturer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ExonumID      string    `json:"exonum_id"`
	DeviceGroupID string    `json:"device_group_id"`
	OwnerKey      string    `json:"owner_key"`
}

type DeviceResponseData struct {
	Devices      []Device      `json:"devices"`
	DeviceGroups []DeviceGroup `json:"device_groups"`
}

type DeviceResponse struct {
	BaseResponse
	Data []DeviceResponseData `json:"data"`
}
