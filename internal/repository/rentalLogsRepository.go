package repository

import (
	"context"
	"milestone2/internal/entity"

	"gorm.io/gorm"
)

type RentalRepo struct {
	db *gorm.DB
}

func NewRentalRepository(db *gorm.DB) *RentalRepo {
	return &RentalRepo{db}
}

func (rr *RentalRepo) GetAll() (logs []entity.RentalLogs, err error) { 
	if err := rr.db.WithContext(context.Background()).Find(&logs).Error; err != nil {
		return []entity.RentalLogs{}, err
	}

	return logs, nil
}