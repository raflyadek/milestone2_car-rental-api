package service

import (
	"fmt"
	"log"
	"milestone2/internal/entity"
	"time"
)

type PaymentRepository interface {
	Create(payment *entity.Payments) (err error)
	GetAll() (payments []entity.Payments, err error)
	GetByUserId(userId int) (payment []entity.Payments, err error)
	GetById(id int) (payment entity.Payments, err error)
	TransactionUpdate(paymentId, totalDay int, availabilityUntil string) (err error)
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
	// check if the car is avail
	getCarInfo, err := ps.carRepo.GetById(req.CarId)
	if err != nil {
		log.Print(err.Error())
		return 
	}

	// if !getCarInfo.Availability {
	// 	return entity.PaymentInfoResponse{}, fmt.Errorf("error %s", err)
	// }

	//check the car availability until
	availUntil := getCarInfo.AvailabilityUntil
	availUntilParsed, err := time.Parse(time.RFC3339, availUntil)
	if err != nil {
		log.Print("avail until")
		return
	}
	//if availuntilparsed after the time is now so if availuntilparse is 20 and now is 21 
	//the value then false and continue and if true then it stops right here.
	availUntilBool := availUntilParsed.After(time.Now())
	if availUntilBool {
		log.Print("check avail until")
		return 
	}



	// price is flexible according to day
	//still error cannot parse 
	//ERRO[0003] parsing time "2025-11-11" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "T" 
	totalDay, err := ps.totalDay(req.EndDate, req.StartDate)
	if err != nil {
		log.Print("here")
		return 
	}

	//valid until
	templateDate := "2006-01-02 15:04:05"
	validUntil := time.Now().Add(time.Minute * 10).Format(templateDate)

	totalPrice := getCarInfo.Price * float64(totalDay)

	payment := entity.Payments{
		UserId: userId,
		CarId: req.CarId,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
		Price: totalPrice,
		ValidUntil: validUntil,
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

func (ps *PaymentServ) GetAllPayment() (resp []entity.PaymentInfoResponse, err error) {
	payments, err := ps.paymentRepo.GetAll()
	if err != nil {
		return []entity.PaymentInfoResponse{}, err
	}

	for _, info := range payments {
		resp = append(resp, entity.PaymentInfoResponse{
			Id: info.Id,
			UserId: info.UserId,
			User: info.User,
			CarId: info.CarId,
			Car: info.Car,
			StartDate: info.StartDate,
			EndDate: info.EndDate,
			Price: info.Price,
			Status: info.Status,
			ValidUntil: info.ValidUntil,
			CreatedAt: info.CreatedAt,
		})
	}
	return resp, nil
}

func (ps *PaymentServ) GetByUserIdPayment(userId int) (resp []entity.PaymentInfoResponse, err error) {
	payments, err := ps.paymentRepo.GetByUserId(userId)
	if err != nil {
		return []entity.PaymentInfoResponse{}, err
	}

	for _, info := range payments {
		resp = append(resp, entity.PaymentInfoResponse{
			Id: info.Id,
			UserId: info.UserId,
			User: info.User,
			CarId: info.CarId,
			Car: info.Car,
			StartDate: info.StartDate,
			EndDate: info.EndDate,
			Price: info.Price,
			Status: info.Status,
			ValidUntil: info.ValidUntil,
			CreatedAt: info.CreatedAt,
		})
	}

	return resp, nil
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
		StartDate: payment.StartDate,
		EndDate: payment.EndDate,
		Price: payment.Price,
		Status: payment.Status,
		CreatedAt: payment.CreatedAt,
		ValidUntil: payment.ValidUntil,
	}
	return getByIdResp, nil
}

func (ps *PaymentServ) TransactionUpdatePayment(paymentId int) (resp entity.PaidPaymentResponse, err error) {
	paymentInfo, err := ps.paymentRepo.GetById(paymentId)
	if err != nil {
		return entity.PaidPaymentResponse{}, err
	}

	//check avail from car.availability OR car.availability.until?? <-
	//we can check with if car.availability.until.day < now then continue
	//yeah using availability.until.day so if someone booked it for next week even 
	//even tho the availability is false anyone can still rented it if < availability.until.day
	// fmt.Printf("data %+v", paymentInfo)
	if !paymentInfo.Car.Availability {
		return entity.PaidPaymentResponse{}, fmt.Errorf("already booked")
	}

	//parsing date
	//formatTime := "2006-01-02 15:04:05"
	formatDate := "2006-01-02"
	endDate := paymentInfo.EndDate
	startDate := paymentInfo.StartDate

	parseEndDate, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		log.Print("error on this parseenddate")
		return
	}

	parseStartDate, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		log.Print("eror start date parse")
		return
	}
	//total day 	
	totalDay := parseEndDate.Day() - parseStartDate.Day()
	//avail until cars + 1 day for the buffer time
	carsAvailUntil := parseEndDate.Add(time.Hour * 24).Format(formatDate)

	//valid until check 
	// formatTime := "15:04:05"
	validUntil := paymentInfo.ValidUntil

	parseValidUntil, err := time.Parse(time.RFC3339, validUntil)
	if err != nil {
		return
	}

	now := time.Now()
	
	if parseValidUntil.After(now) {
		return entity.PaidPaymentResponse{}, fmt.Errorf("expired payment")
	}

	
	if err := ps.paymentRepo.TransactionUpdate(paymentInfo.Id, totalDay, carsAvailUntil); err != nil {
		log.Print("error on this?")
		return entity.PaidPaymentResponse{}, err
	}

	TransactionUpdatePaymentResp := entity.PaidPaymentResponse{
		Id: paymentInfo.Id,
		UserId: paymentInfo.UserId,
		// User: paymentInfo.User,
		CarId: paymentInfo.CarId,
		// Car: paymentInfo.Car,
		TotalDay: totalDay,
		TotalSpent: paymentInfo.Price,
		CreatedAt: paymentInfo.CreatedAt,
	}

	return TransactionUpdatePaymentResp, nil
}