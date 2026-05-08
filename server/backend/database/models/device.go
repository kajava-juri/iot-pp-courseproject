package models

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	DeviceName  string   `gorm:"uniqueIndex;size:100;type:text" json:"device_name" db:"device_name"`
	Name        string   `json:"name" db:"name"`
	UptimeMs    float64  `json:"uptime_ms" db:"uptime_ms"`
	Description string   `json:"description" db:"description"`
	RoomID      *uint    `db:"room_id" json:"room_id,omitempty"`
	Room        *Room    `gorm:"foreignKey:RoomID;references:ID" json:"room"`
	PatientID   *uint    `db:"patient_id" json:"patient_id,omitempty"`
	Patient     *Patient `gorm:"foreignKey:PatientID;references:ID" json:"patient"`
}
