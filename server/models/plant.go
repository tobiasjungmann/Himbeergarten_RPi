package models

import "time"

type Plant struct {
	Plant      int32     `gorm:"primary_key;AUTO_INCREMENT;column:plant;type:int;not null;" json:"plant" `
	Name       string    `gorm:"column:name;type:mediumtext;not null;" json:"name" `
	Info       string    `gorm:"column:info;type:mediumtext;not null;" json:"info" `
	Type       string    `gorm:"column:type;type:mediumtext;not null;" json:"type" `
	Humidity   int32     `gorm:"column:humidity;type:int;not null;" json:"humidity" `
	SensorSlot int32     `gorm:"column:sensorslot;type:int;not null;" json:"sensorslot"`
	Watered    time.Time `gorm:"column:watered;type:timestamp;not null;" json:"watered"`
}

// TableName sets the insert table name for this struct type
func (n *Plant) TableName() string {
	return "plant"
}
