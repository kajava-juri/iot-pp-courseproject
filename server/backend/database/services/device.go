package services

import (
	postgres "backend/database"
	"backend/database/models"
	"fmt"

	"gorm.io/gorm"
)

type DeviceService struct{}

var Device = DeviceService{}

func (s DeviceService) Create(device *models.Device) error {
	if err := postgres.DB().Create(device).Error; err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}
	return nil
}

func (s DeviceService) GetByID(id uint) (*models.Device, error) {
	var device models.Device
	result := postgres.DB().Preload("Room").First(&device, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device %d not found: %w", id, result.Error)
		}
		return nil, fmt.Errorf("failed to get device %d: %w", id, result.Error)
	}
	return &device, nil
}

func (s DeviceService) GetByDeviceID(deviceID string) (*models.Device, error) {
	var device models.Device
	result := postgres.DB().Preload("Room").Where("device_id = ?", deviceID).First(&device)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device %s not found: %w", deviceID, result.Error)
		}
		return nil, fmt.Errorf("failed to get device %s: %w", deviceID, result.Error)
	}
	return &device, nil
}

func (s DeviceService) Delete(id uint) error {
	if err := postgres.DB().Delete(&models.Device{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete device %d: %w", id, err)
	}
	return nil
}
