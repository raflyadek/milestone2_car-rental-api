package service

import (
	"milestone2/internal/entity"

	"github.com/sirupsen/logrus"
)

type RentalLogsRepository interface {
	GetAll() (logs []entity.RentalLogs, err error)
	GetById(id int) (logs entity.RentalLogs, err error)
	GetByUserId(userId int) (logs []entity.RentalLogs, err error)
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
