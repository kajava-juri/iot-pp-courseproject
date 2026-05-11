package services

import (
	postgres "backend/database"
	"backend/database/models"
	"fmt"
	"strings"

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

func (s AlertService) GetActiveByPatientID(patientID uint) ([]models.Alert, error) {
	var alerts []models.Alert
	result := postgres.DB().Where("patient_id = ? AND resolved = ? AND declined = ?", patientID, false, false).
		Preload("Patient").
		Preload("Event").
		Preload("Event.Room").
		Find(&alerts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active alerts for patient %d: %w", patientID, result.Error)
	}
	return alerts, nil
}

func (s AlertService) GetPatientStatus(patientID uint) (string, error) {
	alerts, err := s.GetActiveByPatientID(patientID)
	if err != nil {
		return "", err
	}
	if len(alerts) == 0 {
		return "", nil
	}

	status := alerts[0].Severity
	statusWeight := severityWeight(status)
	for _, alert := range alerts[1:] {
		if weight := severityWeight(alert.Severity); weight > statusWeight {
			status = alert.Severity
			statusWeight = weight
		}
	}

	return status, nil
}

func severityWeight(severity string) int {
	switch strings.ToLower(severity) {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
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
