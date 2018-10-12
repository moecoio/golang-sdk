package db

type Transaction struct {
	ID         int    `json:"id"`
	Hash       string `json:"hash"`
	DeviceHash string `json:"device_hash"`
	Timestamp  int    `json:"timestamp"`
	Uplink     int    `json:"uplink"`
	Sended     int    `json:"sended"`
	Payload    string `json:"payload"`
}

type Device struct {
	ID            string `json:"id"`
	Hash          string `json:"hash"`
	Manufacturer  string `json:"manufacturer"`
	CreatedAt     int    `json:"created_at"`
	UpdatedAt     int    `json:"updated_at"`
	ExonumID      string `json:"exonum_id"`
	DeviceGroupID string `json:"device_group_id"`
	OwnerKey      string `json:"owner_key"`
}

type DeviceGroup struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	GroupType        int    `json:"group_type"`
	UplinkLifetime   int    `json:"uplink_lifetime"`
	DownlinkLifetime int    `json:"downlink_lifetime"`
	Services         string `json:"services"`
	CreatedAt        int    `json:"created_at"`
	UpdatedAt        int    `json:"updated_at"`
	ExonumID         string `json:"exonum_id"`
	OwnerKey         string `json:"owner_key"`
}
