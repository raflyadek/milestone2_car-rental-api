package repository

import (
	"context"
	"milestone2/internal/entity"

	"gorm.io/gorm"
)

type RentalRepo struct {
	db *gorm.DB
}

func NewRentalLogsRepository(db *gorm.DB) *RentalRepo {
	return &RentalRepo{db}
}

func (rr *RentalRepo) GetAll() (logs []entity.RentalLogs, err error) { 
	if err := rr.db.WithContext(context.Background()).
	Omit("User", "Car").
	Find(&logs).Error; err != nil {
		return []entity.RentalLogs{}, err
	}

	return logs, nil
}

func (rr *RentalRepo) GetById(id int) (logs entity.RentalLogs, err error) {
	if err := rr.db.WithContext(context.Background()).
	Omit("User", "Car").
	First(&logs, "id = ?", id).Error; err != nil {
		return entity.RentalLogs{}, err
	}

	return logs, nil
}

func (rr *RentalRepo) GetByUserId(userId int) (logs []entity.RentalLogs, err error) {
	if err := rr.db.WithContext(context.Background()).
	Preload("Car").
	Preload("Car.Categories").
	Find(&logs, "user_id = ?", userId).Error; err != nil {
		return []entity.RentalLogs{}, err
	}

	return logs, nil
} 