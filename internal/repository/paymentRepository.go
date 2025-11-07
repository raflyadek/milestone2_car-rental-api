package repository

import (
	"context"
	"milestone2/internal/entity"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db}
}

func (pr *PaymentRepo) Create(payment *entity.Payments) (err error) {
	if err := pr.db.WithContext(context.Background()).Omit("Status", "CreatedAt").
	Create(payment).Error; err != nil {
		return err
	}

	return nil
}

func (pr *PaymentRepo) GetAll() (payments []entity.Payments, err error) {
	if err := pr.db.WithContext(context.Background()).
	Preload("User").
	Preload("Car").
	Preload("Car.Categories").
	Find(&payments).Error; err != nil {
		return []entity.Payments{}, err
	}
	return payments, nil
}

func (pr *PaymentRepo) GetByUserId(userId int) (payment []entity.Payments, err error) {
	if err := pr.db.WithContext(context.Background()).
	Preload("User").
	Preload("Car").
	Preload("Car.Categories").
	Find(&payment, "user_id = ?", userId).Error; err != nil {
		return []entity.Payments{}, err
	}

	return payment, nil
}

func (pr *PaymentRepo) GetById(id int) (payment entity.Payments, err error) {
	if err := pr.db.WithContext(context.Background()).
	Preload("User").
	Preload("Car").
	Preload("Car.Categories").
	First(&payment, "id = ?", id).Error; err != nil {
		return entity.Payments{}, err
	}

	return payment, nil
}

func (pr *PaymentRepo) TransactionUpdate(paymentId, totalDay int, availabilityUntil string) (err error) {
	errr := pr.db.WithContext(context.Background()).
	Transaction(func(tx *gorm.DB) error {
		//update payment
		var payment entity.Payments
		if err := tx.Model(&payment).Where("id = ?", paymentId).Update("Status", true).Error; err != nil {
			return err
		}

		//get payment info
		if err := tx.First(&payment, "id = ?", paymentId).Error; err != nil {
			return err
		}

		//create rental log
		rentalLog := entity.RentalLogs{
			UserId: payment.UserId,
			CarId: payment.CarId,
			PaymentId: payment.Id,
			StartDate: payment.StartDate,
			EndDate: payment.EndDate,
			TotalDay: totalDay,
			TotalSpent: payment.Price,
		}

		if err := tx.Omit("CreatedAt").Create(&rentalLog).Error; err != nil {
			return err
		}

		//update the car availability until 
		var car entity.Cars
		if err := tx.Model(&car).Where("id = ?", payment.CarId).
		Updates(map[string]interface{}{
			// "Availability": false,
			"AvailabilityUntil": availabilityUntil,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if errr != nil {
		return err
	}

	return nil
}