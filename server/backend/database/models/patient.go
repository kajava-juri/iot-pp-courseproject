package models

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Name      string `gorm:"not null" json:"name"`
	PatientID string `gorm:"uniqueIndex;not null" json:"patient_id"`
	RoomID    uint   `db:"room_id" json:"room_id"`
	Room      Room   `gorm:"foreignKey:RoomID;references:ID" json:"room"`
}
