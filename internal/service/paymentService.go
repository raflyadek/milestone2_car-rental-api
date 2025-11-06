package service

import (
	"fmt"
	"log"
	"milestone2/internal/entity"
	"time"
)

type PaymentRepository interface {
	Create(payment *entity.Payments) (err error)
	GetById(id int) (payment entity.Payments, err error)
	TransactionUpdate(paymentId, totalDay int, availabilityUntil string) (payment entity.Payments, err error)
}

type CarRepository interface {
	GetById(id int) (car entity.Cars, err error)
}

type PaymentServ struct {
	paymentRepo PaymentRepository
	carRepo CarRepository
}

func NewPaymentService(paymentRepo PaymentRepository, carRepo CarRepository) *PaymentServ {
	return &PaymentServ{paymentRepo, carRepo}
}

func (ps *PaymentServ) CreatePayment(userId int, req entity.CreatePaymentRequest) (resp entity.PaymentInfoResponse, err error) {
	//check if the car is avail
	getCarInfo, err := ps.carRepo.GetById(req.CarId)
	if err != nil {
		log.Print(err.Error())
		return 
	}

	if !getCarInfo.Availability {
		return entity.PaymentInfoResponse{}, fmt.Errorf("error %s", err)
	}

	// price is flexible according to day
	totalDay, err := ps.totalDay(req.EndDate, req.StartDate)
	if err != nil {
		log.Print(err.Error())
		return 
	}

	totalPrice := getCarInfo.Price * float64(totalDay)

	payment := entity.Payments{
		UserId: userId,
		CarId: req.CarId,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
		Price: totalPrice,
	}
	if err := ps.paymentRepo.Create(&payment); err != nil {
		return entity.PaymentInfoResponse{}, err
	}

	paymentInfo, err := ps.GetByIdPayment(payment.Id)
	if err != nil {
		return entity.PaymentInfoResponse{}, err
	}

	return paymentInfo, nil
}

func (ps *PaymentServ) GetByIdPayment(id int) (resp entity.PaymentInfoResponse, err error) {
	payment, err := ps.paymentRepo.GetById(id)
	if err != nil {
		return entity.PaymentInfoResponse{}, err
	}

	getByIdResp := entity.PaymentInfoResponse{
		Id: payment.Id,
		UserId: payment.UserId,
		User: payment.User,
		CarId: payment.CarId,
		Car: payment.Car,
		RentalPeriod: payment.RentalPeriod,
		StartDate: payment.StartDate,
		EndDate: payment.EndDate,
		Price: payment.Price,
		Status: payment.Status,
		CreatedAt: payment.CreatedAt,
	}
	return getByIdResp, nil
}

func (ps *PaymentServ) TransactionUpdatePayment(paymentId int) (resp entity.PaidPaymentResponse, err error) {
	paymentInfo, err := ps.paymentRepo.GetById(paymentId)
	if err != nil {
		return entity.PaidPaymentResponse{}, err
	}

	totalDay, err := ps.totalDay(paymentInfo.EndDate, paymentInfo.StartDate)
	if err != nil {
		log.Print(err.Error())
		return
	}

	//avail until cars
	formatTime := "2006-01-02 15:04:05"
	formatDate := "2006-01-02"
	endDate := paymentInfo.EndDate

	parseEndDate, err := time.Parse(formatTime, endDate)
	if err != nil {
		log.Print(err.Error())
		return
	}
	carsAvailUntil := parseEndDate.Add(time.Hour * 24).Format(formatDate)

	TransactionResp, err := ps.paymentRepo.TransactionUpdate(paymentInfo.Id, totalDay, carsAvailUntil)
	if err != nil {
		log.Print(err.Error())
		return 
	}

	TransactionUpdatePaymentResp := entity.PaidPaymentResponse{
		Id: TransactionResp.Id,
		UserId: TransactionResp.UserId,
		User: TransactionResp.User,
		CarId: TransactionResp.CarId,
		Car: TransactionResp.Car,
		TotalDay: totalDay,
		TotalSpent: TransactionResp.Price,
		CreatedAt: TransactionResp.CreatedAt,
	}

	return TransactionUpdatePaymentResp, nil
}