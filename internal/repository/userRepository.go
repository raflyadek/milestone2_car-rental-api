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

func (ur *UserRepo) Create(user entity.User) (register entity.UserRegister, err error) {
	err = ur.db.WithContext(context.Background()).Create(&user).Error
	return register, nil
}

func (ur *UserRepo) GetByEmail(email string) (user entity.User, err error) {
	err = ur.db.WithContext(context.Background()).First(&user, "email = ?", email).Error
	return user, nil
}