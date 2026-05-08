package models

import (
	"time"

	"gorm.io/gorm"
)

// Alert represents an actionable notification derived from events or rules.
// For the demo we keep it lightweight: common fields, optional link to the originating Event.
type Alert struct {
	gorm.Model
	EventID      *uint      `gorm:"index" json:"event_id"`
	Event        *Event     `gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"event,omitempty"`
	PatientID    *uint      `gorm:"index" json:"patient_id"`
	Severity     string     `gorm:"size:50" json:"severity"` // e.g. "low", "medium", "high"
	Message      string     `gorm:"type:text" json:"message"`
	// Below fields are assumed for single staff member for demo
	// In real system, a joint table AlertStaff would be needed for many-to-many relationship
	Acknowledged bool       `gorm:"default:false" json:"acknowledged"`
	Resolved     bool       `gorm:"default:false" json:"resolved"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	Declined	 bool       `gorm:"default:false" json:"declined"`
}
