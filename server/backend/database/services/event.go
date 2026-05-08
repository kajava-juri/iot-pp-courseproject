package services

import (
	postgres "backend/database"
	"backend/database/models"
	"fmt"

	"gorm.io/gorm"
)

type EventService struct{}

var Event = EventService{}

func (s EventService) GetAll() ([]models.Event, error) {
	var events []models.Event
	result := postgres.DB().Preload("Rooms").Preload("Patients").Find(&events)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get events: %w", result.Error)
	}
	return events, nil
}

func (s EventService) Create(event *models.Event) error {
	if err := postgres.DB().Create(event).Error; err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (s EventService) GetByID(id uint) (*models.Event, error) {
	var event models.Event
	result := postgres.DB().Preload("Rooms").Preload("Patients").First(&event, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event %d not found: %w", id, result.Error)
		}
		return nil, fmt.Errorf("failed to get event %d: %w", id, result.Error)
	}
	return &event, nil
}

func (s EventService) Delete(id uint) error {
	if err := postgres.DB().Delete(&models.Event{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete event %d: %w", id, err)
	}
	return nil
}

func (s EventService) List(page int, pageSize int) ([]models.Event, error) {
	// Validate and constrain pageSize
	if pageSize < 50 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	// Ensure page is at least 1
	if page < 1 {
		page = 1
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	var events []models.Event
	result := postgres.DB().Offset(offset).Limit(pageSize).Preload("Rooms").Preload("Patients").Find(&events)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get events: %w", result.Error)
	}
	return events, nil
}
