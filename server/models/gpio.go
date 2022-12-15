package models

type Gpio struct {
	Gpio   int32 `gorm:"primary_key;AUTO_INCREMENT;column:gpio;type:int;not null;" json:"gpio" `
	Device int32 `gorm:"foreignKey:Device;column:device;type:int;not null;" json:"device" `
	Used   bool  `gorm:"foreignKey:Used;column:device;type:bool;not null;" json:"Used" `
}

// TableName sets the insert table name for this struct type
func (n *Gpio) TableName() string {
	return "gpio"
}
