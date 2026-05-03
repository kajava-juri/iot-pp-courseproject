package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	RoomID   string `gorm:"uniqueIndex;not null" json:"room_id"`
	RoomName string `gorm:"not null" json:"room_name"`
}