package service

import (
	"fmt"
	"milestone2/internal/entity"
	"time"

	"github.com/sirupsen/logrus"
)

type RentalLogsRepository interface {
	GetAll() (logs []entity.RentalLogs, err error)
	GetById(id int) (logs entity.RentalLogs, err error)
	GetByUserId(userId int) (logs []entity.RentalLogs, err error)
	GetByCarId(carId int) (logs []entity.RentalLogs, err error)
}

type RentalServ struct {
	logsRepo RentalLogsRepository
}

func NewRentalLogsService(logsRepo RentalLogsRepository) *RentalServ {
	return &RentalServ{logsRepo}
}

func (rs *RentalServ) GetAllLogs() (resp []entity.RentalLogsResponseAdmin, err error) {
	logs, err := rs.logsRepo.GetAll()
	if err != nil {
		logrus.Printf("error get all logs %s", err)
		return []entity.RentalLogsResponseAdmin{}, err
	}

	for _, log := range logs {
		resp = append(resp, entity.RentalLogsResponseAdmin{
			Id: log.Id,
			UserId: log.UserId,
			CarId: log.CarId,
			PaymentId: log.PaymentId,
			StartDate: log.StartDate,
			EndDate: log.EndDate,
			TotalDay: log.TotalDay,
			TotalSpent: log.TotalSpent,
			CreatedAt: log.CreatedAt,
		})
	}

	return resp, nil
}

func (rs *RentalServ) GetByIdLogs(id int) (resp entity.RentalLogsResponseAdmin, err error) {
	log, err := rs.logsRepo.GetById(id)
	if err != nil {
		logrus.Printf("error get by id logs %s", err)
		return entity.RentalLogsResponseAdmin{}, err
	}

	logResp := entity.RentalLogsResponseAdmin{
		Id: log.Id,
		UserId: log.UserId,
		CarId: log.CarId,
		PaymentId: log.PaymentId,
		StartDate: log.StartDate,
		EndDate: log.EndDate,
		TotalDay: log.TotalDay,
		TotalSpent: log.TotalSpent,
		CreatedAt: log.CreatedAt,
	}

	return logResp, nil
}

func (rs *RentalServ) GetByUserIdLogs(userId int) (resp []entity.RentalLogsResponseUser, err error) {
	logs, err := rs.logsRepo.GetByUserId(userId)
	if err != nil {
		logrus.Printf("error get logs by user id %s", err)
		return []entity.RentalLogsResponseUser{}, err
	}

	for _, log := range logs {
		resp = append(resp, entity.RentalLogsResponseUser(log))
	}

	return resp, nil
}


func (rs *RentalServ) CheckAvailabilityByCarId(req entity.CheckCarAvailabilityRequest) (resp entity.RentalAvailabilityResponse, err error) {
	logs, err := rs.logsRepo.GetByCarId(req.CarId)
	if err != nil {
		logrus.Print("error get log by")
	}
	layout := "2006-01-02"
	startDate := req.StartDate
	endDate := req.EndDate

	startDateParse, errr := time.Parse(layout, startDate)
	if errr != nil {
		logrus.Printf("failed parse start")
		return entity.RentalAvailabilityResponse{}, err
	}

	endDateParse, errrr := time.Parse(layout, endDate)
	if errrr != nil {
		logrus.Printf("failed parse end")
		return entity.RentalAvailabilityResponse{}, errr
	}

	

	for _, log := range logs {
		endDateLogParse, err := time.Parse(time.RFC3339, log.EndDate)
		startDateLogParse, errr := time.Parse(time.RFC3339, log.StartDate)
		fmt.Print(endDateLogParse, startDateLogParse, startDateParse, endDateParse)
		if err != nil {
			logrus.Printf("failed parse end log")
			return entity.RentalAvailabilityResponse{}, err
		}

		if errr != nil {
			logrus.Printf("failed parse start log")
			return entity.RentalAvailabilityResponse{}, errr
		}
        // check overlap
        if !startDateParse.Before(endDateLogParse) && endDateParse.After(startDateLogParse) {
            resp.Availability = true
            return resp, nil
        }	
    }

	return resp, nil
}