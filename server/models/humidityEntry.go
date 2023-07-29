package models

import "time"

type HumidityEntry struct {
	HumidityEntry int32     `gorm:"primary_key;AUTO_INCREMENT;column:humidityEntry;type:int;not null;" json:"humidityEntry" `
	DeviceID      int32     `gorm:"column:device;type:int;not null;" json:"device" `
	SensorSlot    int32     `gorm:"column:sensor;type:int;not null;" json:"sensor" `
	Value         int32     `gorm:"column:value;type:int;not null;" json:"value" `
	Timestamp     time.Time `gorm:"column:timestamp;type:timestamp;not null;" json:"timestamp"`
}

// TableName sets the insert table name for this struct type
func (n *HumidityEntry) TableName() string {
	return "humidity_entry"
}
