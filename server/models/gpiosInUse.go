package models

type GpioInUse struct {
	GpioInUse int32  `gorm:"primary_key;AUTO_INCREMENT;column:gpioInUse;type:int;not null;" json:"gpioInUse" `
	Plant     int32  `gorm:"foreignKey:Plant;column:plant;type:int;not null;" json:"plant" `
	Gpio      int32  `gorm:"foreignKey:Gpio;column:gpio;type:int;not null;" json:"gpio" `
	Interface string `gorm:"foreignKey:Gpio;column:interface;type:mediumtext;not null;" json:"interface" `
}

// TableName sets the insert table name for this struct type
func (n *GpioInUse) TableName() string {
	return "gpio_in_use"
}
