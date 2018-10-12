package types

import (
	"clients/prot"
	"db"
	"encoding/json"
	"time"
)

func timeToInt(t time.Time) int {
	return int(t.Unix())
}

func intToTime(i int) time.Time {
	return time.Unix(int64(i), 0)
}

func intToBool(i int) bool {
	return i != 0
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func TransactionToReq(t db.Transaction) prot.TransactionReq {
	return prot.TransactionReq{
		ID:         t.ID,
		//Hash:       t.Hash,
		DeviceHash: t.DeviceHash,
		Timestamp:  intToTime(t.Timestamp),
		Uplink:     intToBool(t.Uplink),
		Payload:    t.Payload,
		Status:     1,
	}
}

func TrasactionsToReq(transactions []db.Transaction) []prot.TransactionReq {
	res := make([]prot.TransactionReq, 0, len(transactions))
	for _, v := range transactions {
		res = append(res, TransactionToReq(v))
	}
	return res
}

func DeviceFromResponse(device prot.Device) db.Device {
	return db.Device{
		ID:            device.ID,
		Hash:          device.Hash,
		Manufacturer:  device.Manufacturer,
		CreatedAt:     timeToInt(device.CreatedAt),
		UpdatedAt:     timeToInt(device.UpdatedAt),
		ExonumID:      device.ExonumID,
		DeviceGroupID: device.DeviceGroupID,
		OwnerKey:      device.OwnerKey,
	}
}

func DevicesFromResponse(devices []prot.Device) []db.Device {
	ret := make([]db.Device, 0, len(devices))
	for _, v := range devices {
		ret = append(ret, DeviceFromResponse(v))
	}
	return ret
}

func DeviceToResponse(device db.Device) prot.Device {
	return prot.Device{
		ID:            device.ID,
		Hash:          device.Hash,
		Manufacturer:  device.Manufacturer,
		CreatedAt:     intToTime(device.CreatedAt),
		UpdatedAt:     intToTime(device.UpdatedAt),
		ExonumID:      device.ExonumID,
		DeviceGroupID: device.DeviceGroupID,
		OwnerKey:      device.OwnerKey,
	}
}

func DevicesToResponse(devices []db.Device) []prot.Device {
	ret := make([]prot.Device, 0, len(devices))
	for _, v := range devices {
		ret = append(ret, DeviceToResponse(v))
	}
	return ret
}

func DeviceGroupFromResponse(deviceGroup prot.DeviceGroup) (*db.DeviceGroup, error) {
	serviceStr, err := json.Marshal(deviceGroup.Services)
	if err != nil {
		return nil, err
	}
	return &db.DeviceGroup{
		ID:               deviceGroup.ID,
		Name:             deviceGroup.Name,
		GroupType:        deviceGroup.GroupType,
		UplinkLifetime:   deviceGroup.UplinkLifetime,
		DownlinkLifetime: deviceGroup.DownlinkLifetime,
		Services: string(serviceStr),
		CreatedAt: timeToInt(deviceGroup.CreatedAt),
		UpdatedAt: timeToInt(deviceGroup.UpdatedAt),
		ExonumID: deviceGroup.ExonumID,
		OwnerKey: deviceGroup.OwnerKey,
	}, nil
}

func DeviceGroupsFromResponse(deviceGroups []prot.DeviceGroup) ([]db.DeviceGroup, error)  {
	ret := make([]db.DeviceGroup, 0, len(deviceGroups))
	for _, v := range deviceGroups {
		tmp, err :=  DeviceGroupFromResponse(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *tmp)
	}
	return ret, nil
}

func DeviceGroupToResponse(deviceGroup db.DeviceGroup) (*prot.DeviceGroup, error) {
	var services []prot.Service
	err := json.Unmarshal([]byte(deviceGroup.Services), &services)
	if err != nil {
		return nil, err
	}
	return &prot.DeviceGroup{
		ID: deviceGroup.ID,
		Name: deviceGroup.Name,
		GroupType: deviceGroup.GroupType,
		UplinkLifetime: deviceGroup.UplinkLifetime,
		DownlinkLifetime: deviceGroup.DownlinkLifetime,
		Services: services,
		CreatedAt: intToTime(deviceGroup.CreatedAt),
		UpdatedAt: intToTime(deviceGroup.UpdatedAt),
		ExonumID: deviceGroup.ExonumID,
		OwnerKey: deviceGroup.OwnerKey,
	}, nil
}

func DeviceGroupsToResponse(deviceGroups []db.DeviceGroup) ([]prot.DeviceGroup, error)  {
	ret := make([]prot.DeviceGroup, 0, len(deviceGroups))
	for _, v := range deviceGroups {
		tmp, err :=  DeviceGroupToResponse(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *tmp)
	}
	return ret, nil
}
