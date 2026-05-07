package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	RoomName string `gorm:"not null" json:"room_name" db:"room_name"`
}
