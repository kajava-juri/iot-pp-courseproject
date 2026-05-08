package models

import (
	"gorm.io/gorm"
)

// EventType gives the event kind a stable type instead of a plain string.
type EventType string

const (
	EventTypeFall                EventType = "fall"
	EventTypeAbnormalTemperature EventType = "abnormal_temperature"
	EventTypeTemperatureReading  EventType = "temperature_reading"
	EventTypeMotionDetected      EventType = "motion_detected"
)

// Event represents a raw device or system event.
// For the demo we keep a single optional room and patient relation.
type Event struct {
	gorm.Model
	DeviceID  uint      `gorm:"index" json:"device_id" db:"device_id"`
	Device    Device    `gorm:"foreignKey:DeviceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"device"`
	Type      EventType `gorm:"type:text;index;not null" json:"type"`
	RoomID    *uint     `gorm:"index" json:"room_id,omitempty" db:"room_id"`
	Room      *Room     `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"room,omitempty"`
	PatientID *uint     `gorm:"index" json:"patient_id,omitempty" db:"patient_id"`
	Patient   *Patient  `gorm:"foreignKey:PatientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"patient,omitempty"`
}
