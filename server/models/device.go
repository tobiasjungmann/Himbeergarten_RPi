package models

import "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"

type Device struct {
	Device int32             `gorm:"primary_key;AUTO_INCREMENT;column:device;type:int;not null;" json:"device" `
	Type   proto.DeviceTypes `gorm:"column:type;type:enum('DeviceTypes_DEVICE_UNDEFINED','DeviceTypes_DEVICE_RPI','DeviceTypes_DEVICE_ARDUINO_NANO');not null;" json:"type" `
}

// TableName sets the insert table name for this struct type
func (n *Device) TableName() string {
	return "device"
}
