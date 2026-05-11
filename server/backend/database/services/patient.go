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

func (s PatientService) List() ([]models.Patient, error) {
	var patients []models.Patient
	if err := postgres.DB().Preload("Room").Find(&patients).Error; err != nil {
		return nil, fmt.Errorf("failed to list patients: %w", err)
	}
	return patients, nil
}

func (s PatientService) GetByHealthID(healthID string) (*models.Patient, error) {
	var patient models.Patient
	result := postgres.DB().Preload("Room").Where("health_id = ?", healthID).First(&patient)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("patient with health_id %s not found: %w", healthID, result.Error)
		}
		return nil, fmt.Errorf("failed to get patient with health_id %s: %w", healthID, result.Error)
	}
	return &patient, nil
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

func (s PatientService) UpdateStatus(id uint, status string) (*models.Patient, bool, error) {
	patient, err := s.GetByID(id)
	if err != nil {
		return nil, false, err
	}

	if patient.Status == status {
		return patient, false, nil
	}

	if err := postgres.DB().Model(&models.Patient{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return nil, false, fmt.Errorf("failed to update patient %d status: %w", id, err)
	}

	updatedPatient, err := s.GetByID(id)
	if err != nil {
		return nil, false, err
	}

	return updatedPatient, true, nil
}

func (s PatientService) Delete(id uint) error {
	if err := postgres.DB().Delete(&models.Patient{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete patient %d: %w", id, err)
	}
	return nil
}
