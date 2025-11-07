package service_test

import (
	"errors"
	"milestone2/internal/entity"
	"milestone2/internal/service"
	"milestone2/internal/mocks"
	"testing"
	"time"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayment(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	carRepo := new(mocks.CarRepository)
	svc := service.NewPaymentService(paymentRepo, carRepo)

	car := entity.Cars{
		Id: 1, Price: 100, AvailabilityUntil: time.Now().Add(-time.Hour * 24).Format(time.RFC3339),
	}
	req := entity.CreatePaymentRequest{
		CarId: 1, StartDate: "2025-11-20", EndDate: "2025-11-22",
	}

	carRepo.On("GetById", 1).Return(car, nil)
	paymentRepo.On("Create", mock.Anything).Return(nil)
	paymentRepo.On("GetById", mock.AnythingOfType("int")).Return(entity.Payments{
		Id:        1,
		UserId:    1,
		CarId:     1,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Price:     200,
	}, nil)

	resp, err := svc.CreatePayment(1, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.CarId)
	assert.Equal(t, 200.0, resp.Price)
}

func TestCreatePayment_InvalidCar(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	carRepo := new(mocks.CarRepository)
	svc := service.NewPaymentService(paymentRepo, carRepo)

	carRepo.On("GetById", 1).Return(entity.Cars{}, errors.New("not found"))

	_, err := svc.CreatePayment(1, entity.CreatePaymentRequest{CarId: 1})
	assert.Error(t, err)
}

func TestGetAllPayment(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	carRepo := new(mocks.CarRepository)
	svc := service.NewPaymentService(paymentRepo, carRepo)

	payments := []entity.Payments{{Id: 1, UserId: 2, Price: 500}}
	paymentRepo.On("GetAll").Return(payments, nil)

	resp, err := svc.GetAllPayment()
	assert.NoError(t, err)
	assert.Equal(t, 1, resp[0].Id)
	assert.Equal(t, 500.0, resp[0].Price)
}

func TestGetByUserIdPayment(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	carRepo := new(mocks.CarRepository)
	svc := service.NewPaymentService(paymentRepo, carRepo)

	payments := []entity.Payments{{Id: 1, UserId: 2, Price: 800}}
	paymentRepo.On("GetByUserId", 2).Return(payments, nil)

	resp, err := svc.GetByUserIdPayment(2)
	assert.NoError(t, err)
	assert.Equal(t, 800.0, resp[0].Price)
}

func TestGetByIdPayment(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	carRepo := new(mocks.CarRepository)
	svc := service.NewPaymentService(paymentRepo, carRepo)

	payment := entity.Payments{Id: 1, Price: 700}
	paymentRepo.On("GetById", 1).Return(payment, nil)

	resp, err := svc.GetByIdPayment(1)
	assert.NoError(t, err)
	assert.Equal(t, 700.0, resp.Price)
}

func TestTransactionUpdatePayment(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	carRepo := new(mocks.CarRepository)
	svc := service.NewPaymentService(paymentRepo, carRepo)

	now := time.Now()
	payment := entity.Payments{
		Id: 1, UserId: 1, CarId: 2, Price: 400, ValidUntil: now.Add(-time.Hour).Format(time.RFC3339),
		StartDate: now.Add(-48 * time.Hour).Format(time.RFC3339),
		EndDate:   now.Add(-24 * time.Hour).Format(time.RFC3339),
		Car:       entity.Cars{Availability: true},
	}
	paymentRepo.On("GetById", 1).Return(payment, nil)
	paymentRepo.On("TransactionUpdate", 1, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)

	resp, err := svc.TransactionUpdatePayment(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.UserId)
	assert.Equal(t, 400.0, resp.TotalSpent)
}

func TestTransactionUpdatePayment_AlreadyBooked(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	carRepo := new(mocks.CarRepository)
	svc := service.NewPaymentService(paymentRepo, carRepo)

	payment := entity.Payments{
		Id: 1, Car: entity.Cars{Availability: false},
	}
	paymentRepo.On("GetById", 1).Return(payment, nil)

	_, err := svc.TransactionUpdatePayment(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already booked")
}
