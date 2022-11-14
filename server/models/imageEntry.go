package models

type ImageEntry struct {
	HumidityEntry int32  `gorm:"primary_key;AUTO_INCREMENT;column:humidityEntry;type:int;not null;" json:"humidityEntry" `
	Plant         int32  `gorm:"foreignKey:Plant;column:plant;type:int;not null;" json:"plant" `
	Path          string `gorm:"column:path;type:mediumtext;not null;" json:"path" `
}

// TableName sets the insert table name for this struct type
func (n *ImageEntry) TableName() string {
	return "image_entry"
}
