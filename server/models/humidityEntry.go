package models

import "time"

type HumidityEntry struct {
	HumidityEntry int32     `gorm:"primary_key;AUTO_INCREMENT;column:humidityEntry;type:int;not null;" json:"humidityEntry" `
	Plant         int32     `gorm:"foreignKey:Plant;column:plant;type:int;not null;" json:"plant" `
	Value         int32     `gorm:"column:value;type:int;not null;" json:"value" `
	Timestamp     time.Time `gorm:"column:timestamp;type:timestamp;not null;" json:"timestamp"`
}

// TableName sets the insert table name for this struct type
func (n *HumidityEntry) TableName() string {
	return "humidity_entry"
}
