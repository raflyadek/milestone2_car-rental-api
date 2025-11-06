package service

import (
	"log"
	"milestone2/internal/entity"
	"time"
)

type PaymentRepository interface {
	Create(payment *entity.Payments) (err error)
	TransactionUpdate(paymentId, totalDay int) (payment entity.Payments, err error)
	GetById(id int) (payment entity.Payments, err error)
}

type PaymentServ struct {
	paymentRepo PaymentRepository
}

func NewPaymentService(paymentRepo PaymentRepository) *PaymentServ {
	return &PaymentServ{paymentRepo}
}

func (ps *PaymentServ) CreatePayment(userId int, req entity.CreatePaymentRequest) (resp entity.PaymentInfoResponse, err error) {
	payment := entity.Payments{
		UserId: userId,
		CarId: req.CarId,
		RentalPeriod: req.RentalPeriod,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
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
	endDateDay := paymentInfo.EndDate
	startDateDay := paymentInfo.StartDate
	templateDate := "2006-01-02"

	parseEndDate, err := time.Parse(templateDate, endDateDay)
	if err != nil {
		log.Print(err.Error())
		return entity.PaidPaymentResponse{}, err
	}

	parseStartDate, err := time.Parse(templateDate, startDateDay)
	if err != nil {
		log.Print(err.Error())
	}

	getDayFromEndDate := parseEndDate.Day()
	getDayFromStartDate := parseStartDate.Day()

	totalDay := getDayFromEndDate - getDayFromStartDate

	// totalDay := paymentInfo.EndDate - paymentInfo.StartDate
	TransactionResp, err := ps.paymentRepo.TransactionUpdate(paymentInfo.Id, totalDay)

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