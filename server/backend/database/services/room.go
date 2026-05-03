package services

import (
	postgres "backend/database"
	"backend/database/models"
	"fmt"

	"gorm.io/gorm"
)

type RoomService struct{}

var Room = RoomService{}

func (s RoomService) Create(room *models.Room) error {
	if err := postgres.DB().Create(room).Error; err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

func (s RoomService) GetByID(id uint) (*models.Room, error) {
	var room models.Room
	result := postgres.DB().First(&room, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("room %d not found: %w", id, result.Error)
		}
		return nil, fmt.Errorf("failed to get room %d: %w", id, result.Error)
	}
	return &room, nil
}

func (s RoomService) Delete(id uint) error {
	if err := postgres.DB().Delete(&models.Room{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete room %d: %w", id, err)
	}
	return nil
}
