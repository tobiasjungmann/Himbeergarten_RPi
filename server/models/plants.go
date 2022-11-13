package main

type Plant struct {
	Plant      int32  `gorm:"primary_key;AUTO_INCREMENT;column:plant;type:int;not null;" json:"plant" `
	Name       string `gorm:"column:name;type:mediumtext;not null;" json:"name" `
	Humidity   int32  `gorm:"column:humidity;type:int;not null;" json:"humidity" `
	SensorSlot int32  `gorm:"column:sensorslot;type:int;not null;" json:"sensorslot" `
}

// TableName sets the insert table name for this struct type
func (n *Plant) TableName() string {
	return "plant"
}
