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
	if err := pr.db.WithContext(context.Background()).Omit("Status, CreatedAt").
	Create(payment).Error; err != nil {
		return err
	}

	return nil
}

func (pr *PaymentRepo) GetById(id int) (payment entity.Payments, err error) {
	if err := pr.db.WithContext(context.Background()).First(&payment, "id = ?", id).Error; err != nil {
		return entity.Payments{}, err
	}

	return payment, nil
}

func (pr *PaymentRepo) TransactionUpdate(paymentId, totalDay int) (payment entity.Payments, err error) {
	errr := pr.db.WithContext(context.Background()).
	Transaction(func(tx *gorm.DB) error {
		//update payment
		if err := tx.Model(&payment).Where("id = ?", paymentId).Update("Status", true).Error; err != nil {
			return err
		}

		//get payment info
		var payment entity.Payments
		if err := tx.First(&payment, "id = ?", paymentId).Error; err != nil {
			return err
		}

		//create rental log
		rentalLog := entity.RentalLogs{
			UserId: payment.UserId,
			User: payment.User,
			CarId: payment.CarId,
			Car: payment.Car,
			PaymentId: payment.Id,
			TotalDay: totalDay,
			TotalSpent: payment.Price,
		}

		if err := tx.Create(&rentalLog).Error; err != nil {
			return err
		}

		//update the car avail as false
		var car entity.Cars
		if err := tx.Model(&car).Where("id = ?", payment.CarId).Update("Availability", false).Error; err != nil {
			return err
		}
		return nil
	})
	if errr != nil {
		return entity.Payments{}, err
	}

	return payment, err
}