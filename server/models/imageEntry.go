package models

type ImageEntry struct {
	ImageEntry int32  `gorm:"primary_key;AUTO_INCREMENT;column:imageEntry;type:int;not null;" json:"imageEntry" `
	Plant      int32  `gorm:"foreignKey:Plant;column:plant;type:int;not null;" json:"plant" `
	Path       string `gorm:"column:path;type:mediumtext;not null;" json:"path" `
}

// TableName sets the insert table name for this struct type
func (n *ImageEntry) TableName() string {
	return "image_entry"
}
