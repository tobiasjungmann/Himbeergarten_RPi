package models

type Sensor struct {
	Sensor     int32  `gorm:"primary_key;AUTO_INCREMENT;column:sensor;type:int;not null;" json:"sensor" `
	DeviceId   string `gorm:"column:device;type:int;not null;size:24" json:"device" `
	SensorSlot int32  `gorm:"column:sensor_slot;type:int;not null;" json:"sensor_slot" `
	InUse      bool   `gorm:"column:used;type:bool;not null;" json:"used" `
}

func (n *Sensor) TableName() string {
	return "sensor"
}
