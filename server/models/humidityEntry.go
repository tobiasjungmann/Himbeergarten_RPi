package models

import "time"

type HumidityEntry struct {
	HumidityEntry  int32     `gorm:"primary_key;AUTO_INCREMENT;column:humidityEntry;type:int;not null;" json:"humidityEntry" `
	Sensor         int32     `gorm:"foreignKey:Sensor;column:sensor;type:int;not null;" json:"sensor" `
	Value          int32     `gorm:"column:value;type:int;not null;" json:"value" `
	ValueInPercent int32     `gorm:"column:value_in_percent;type:int;not null;" json:"valueInPercent" `
	Timestamp      time.Time `gorm:"column:timestamp;type:timestamp;not null;" json:"timestamp"`
}

// TableName sets the insert table name for this struct type
func (n *HumidityEntry) TableName() string {
	return "humidity_entry"
}
