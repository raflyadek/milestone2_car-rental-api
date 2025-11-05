package repository

import (
	"context"
	"milestone2/internal/entity"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (ur *UserRepo) Create(user *entity.User) (err error) {
	if err := ur.db.WithContext(context.Background()).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) GetByEmail(email string) (user entity.User, err error) {
	if err := ur.db.WithContext(context.Background()).First(&user, "email = ?", email).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (ur *UserRepo) GetById(id int) (user entity.User, err error) {
	if err := ur.db.WithContext(context.Background()).First(&user, "id = ?", id).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (ur *UserRepo) UpdateValidationStatus(code, email string) (err error) {
	var user entity.User
	err = ur.db.WithContext(context.Background()).Model(&user).
	Where("validation_code = ? AND email = ?", code, email).
	Update("validation_code", true).Error
	if err != nil {
		return err
	}

	return nil
}