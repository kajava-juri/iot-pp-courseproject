package models

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Name        string  `json:"name" db:"name"`
	UptimeMs    float64 `json:"uptime_ms" db:"uptime_ms"`
	Description string  `json:"description" db:"description"`
	RoomID    	uint   `db:"room_id" json:"room_id"`
	Room        Room    `gorm:"foreignKey:RoomID;references:ID" json:"room"`
}
