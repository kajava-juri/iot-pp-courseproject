package services

import (
	postgres "backend/database"
	"backend/database/models"
	"fmt"

	"gorm.io/gorm"
)

type AlertService struct{}

type AlarmService = AlertService

var Alert = AlertService{}

func (s AlertService) GetAll(alert *models.Alert) ([]models.Alert, error) {
	var alerts []models.Alert
	result := postgres.DB().Where(alert, "resolved", "declined").Preload("Patient").Preload("Event").Preload("Event.Room").Find(&alerts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", result.Error)
	}
	return alerts, nil
}

func (s AlertService) Create(alert *models.Alert) error {
	if err := postgres.DB().Create(alert).Error; err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}
	return nil
}

func (s AlertService) Update(alert *models.Alert) error {
	if alert == nil {
		return fmt.Errorf("alert cannot be nil")
	}
	if alert.ID == 0 {
		return fmt.Errorf("alert ID is required")
	}

	// updates := map[string]any{
	// 	"event_id":     alert.EventID,
	// 	"patient_id":   alert.PatientID,
	// 	"severity":     alert.Severity,
	// 	"message":      alert.Message,
	// 	"acknowledged": alert.Acknowledged,
	// 	"resolved":     alert.Resolved,
	// 	"resolved_at":  alert.ResolvedAt,
	// }

	result := postgres.DB().Model(&models.Alert{}).Where("id = ?", alert.ID).Updates(&alert)
	if result.Error != nil {
		return fmt.Errorf("failed to update alert %d: %w", alert.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("alert %d not found", alert.ID)
	}

	return nil
}

func (s AlertService) GetByID(id uint) (*models.Alert, error) {
	var alert models.Alert
	result := postgres.DB().Preload("Patient").Preload("Event").First(&alert, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("alert %d not found: %w", id, result.Error)
		}
		return nil, fmt.Errorf("failed to get alert %d: %w", id, result.Error)
	}
	return &alert, nil
}

func (s AlertService) Delete(id uint) error {
	if err := postgres.DB().Delete(&models.Alert{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete alert %d: %w", id, err)
	}
	return nil
}
