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
// A single event can be linked to multiple rooms and patients through join tables.
type Event struct {
	gorm.Model
	DeviceID string    `gorm:"index;size:100;not null" json:"device_id"`
	Type     EventType `gorm:"type:text;index;not null" json:"type"`
	Rooms    []Room    `gorm:"many2many:event_rooms;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"rooms,omitempty"`
	Patients []Patient `gorm:"many2many:event_patients;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"patients,omitempty"`
}
