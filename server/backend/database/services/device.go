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

func (s DeviceService) List() ([]models.Device, error) {
	var devices []models.Device
	if err := postgres.DB().Preload("Room").Preload("Patient").Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	return devices, nil
}

func (s DeviceService) GetByID(id uint) (*models.Device, error) {
	var device models.Device
	result := postgres.DB().Preload("Room").Preload("Patient").First(&device, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device %d not found: %w", id, result.Error)
		}
		return nil, fmt.Errorf("failed to get device %d: %w", id, result.Error)
	}
	return &device, nil
}

func (s DeviceService) GetByDeviceName(deviceName string) (*models.Device, error) {
	var device models.Device
	result := postgres.DB().Preload("Room").Preload("Patient").Where("device_name = ?", deviceName).First(&device)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device %s not found: %w", deviceName, result.Error)
		}
		return nil, fmt.Errorf("failed to get device %s: %w", deviceName, result.Error)
	}
	return &device, nil
}

func (s DeviceService) Delete(id uint) error {
	if err := postgres.DB().Delete(&models.Device{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete device %d: %w", id, err)
	}
	return nil
}

func (s DeviceService) Update(device *models.Device) error {
	if device == nil {
		return fmt.Errorf("device cannot be nil")
	}
	if device.ID == 0 {
		return fmt.Errorf("device ID is required")
	}

	// Only allow updating mutable fields: Description, RoomID, PatientID
	updates := map[string]any{
		"description": device.Description,
		"room_id":     device.RoomID,
		"patient_id":  device.PatientID,
	}

	result := postgres.DB().Model(&models.Device{}).Where("id = ?", device.ID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update device %d: %w", device.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("device %d not found", device.ID)
	}

	return nil
}
