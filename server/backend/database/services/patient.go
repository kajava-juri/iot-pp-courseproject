package services

import (
	postgres "backend/database"
	"backend/database/models"
	"fmt"

	"gorm.io/gorm"
)

type PatientService struct{}

var Patient = PatientService{}

func (s PatientService) Create(patient *models.Patient) error {
	if err := postgres.DB().Create(patient).Error; err != nil {
		return fmt.Errorf("failed to create patient: %w", err)
	}
	return nil
}

func (s PatientService) GetByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	result := postgres.DB().Preload("Room").First(&patient, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("patient %d not found: %w", id, result.Error)
		}
		return nil, fmt.Errorf("failed to get patient %d: %w", id, result.Error)
	}
	return &patient, nil
}

func (s PatientService) Delete(id uint) error {
	if err := postgres.DB().Delete(&models.Patient{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete patient %d: %w", id, err)
	}
	return nil
}
