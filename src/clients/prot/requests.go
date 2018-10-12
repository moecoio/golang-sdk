package prot

import "time"

type TransactionReq struct {
	ID         int       `json:"id"`
	//Hash       string    `json:"hash"`
	DeviceHash string    `json:"device_hash"`
	Timestamp  time.Time `json:"timestamp"`
	Uplink     bool      `json:"uplink"`
	Payload    string    `json:"payload"`
	Status     int       `json:"status"`
}

type Transactions struct {
	Transactions []TransactionReq `json:"transactions"`
}

type InitGateReq struct {
	APIKey string `json:"api_key"`
	Gate   Gate   `json:"gate"`
}

type Gate struct {
	Hash string `json:"hash"`
}
