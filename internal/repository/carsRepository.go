package repository

import (
	"context"
	"milestone2/internal/entity"

	"gorm.io/gorm"
)

type CarsRepo struct {
	db *gorm.DB
}

func NewCarsRepository(db *gorm.DB) *CarsRepo {
	return &CarsRepo{db}
}

func (cr *CarsRepo) Create(car *entity.Cars) (err error) {
	if err := cr.db.WithContext(context.Background()).Omit("Availability", "AvailabilityUntil").
	Create(car).Error; err != nil {
		return err
	}

	return nil
}

func (cr *CarsRepo) GetAll() (cars []entity.Cars, err error) {
	if err := cr.db.WithContext(context.Background()).Preload("Categories").
	Find(&cars).Error; err != nil {
		return []entity.Cars{}, err
	}

	return cars, nil
}

func (cr *CarsRepo) GetById(id int) (car entity.Cars, err error) {
	if err := cr.db.WithContext(context.Background()).Preload("Categories").
	First(&car, "id = ?", id).Error; err != nil {
		return entity.Cars{}, err
	}

	return car, nil
}