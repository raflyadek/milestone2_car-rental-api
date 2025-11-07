package service_test

import (
	"errors"
	"milestone2/internal/entity"
	"milestone2/internal/mocks"
	"milestone2/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllLogs_Success(t *testing.T) {
	mockRepo := new(mocks.RentalLogsRepository)
	s := service.NewRentalLogsService(mockRepo)

	expectedLogs := []entity.RentalLogs{
		{
			Id:         1,
			UserId:     2,
			CarId:      3,
			PaymentId:  4,
			StartDate:  "2025-11-01",
			EndDate:    "2025-11-03",
			TotalDay:   2,
			TotalSpent: 300000,
			CreatedAt:  "2025-11-01",
		},
	}

	mockRepo.On("GetAll").Return(expectedLogs, nil)

	resp, err := s.GetAllLogs()

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, expectedLogs[0].UserId, resp[0].UserId)
	mockRepo.AssertExpectations(t)
}

func TestGetAllLogs_Error(t *testing.T) {
	mockRepo := new(mocks.RentalLogsRepository)
	s := service.NewRentalLogsService(mockRepo)

	mockRepo.On("GetAll").Return(nil, errors.New("db error"))

	resp, err := s.GetAllLogs()

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
	assert.Len(t, resp, 0)
	mockRepo.AssertExpectations(t)
}

func TestGetByIdLogs_Success(t *testing.T) {
	mockRepo := new(mocks.RentalLogsRepository)
	s := service.NewRentalLogsService(mockRepo)

	expectedLog := entity.RentalLogs{
		Id:         10,
		UserId:     11,
		CarId:      12,
		PaymentId:  13,
		StartDate:  "2025-11-05",
		EndDate:    "2025-11-06",
		TotalDay:   1,
		TotalSpent: 150000,
		CreatedAt:  "2025-11-05",
	}

	mockRepo.On("GetById", expectedLog.Id).Return(expectedLog, nil)

	resp, err := s.GetByIdLogs(expectedLog.Id)

	assert.NoError(t, err)
	assert.Equal(t, expectedLog.PaymentId, resp.PaymentId)
	assert.Equal(t, expectedLog.TotalSpent, resp.TotalSpent)
	mockRepo.AssertExpectations(t)
}

func TestGetByIdLogs_Error(t *testing.T) {
	mockRepo := new(mocks.RentalLogsRepository)
	s := service.NewRentalLogsService(mockRepo)

	mockRepo.On("GetById", 99).Return(entity.RentalLogs{}, errors.New("not found"))

	resp, err := s.GetByIdLogs(99)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.Equal(t, entity.RentalLogsResponseAdmin{}, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetByUserIdLogs_Success(t *testing.T) {
	mockRepo := new(mocks.RentalLogsRepository)
	s := service.NewRentalLogsService(mockRepo)

	expectedLogs := []entity.RentalLogs{
		{Id: 1, UserId: 20, CarId: 2, TotalDay: 3, TotalSpent: 500000},
		{Id: 2, UserId: 20, CarId: 3, TotalDay: 2, TotalSpent: 300000},
	}

	mockRepo.On("GetByUserId", 20).Return(expectedLogs, nil)

	resp, err := s.GetByUserIdLogs(20)

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, expectedLogs[0].CarId, resp[0].CarId)
	mockRepo.AssertExpectations(t)
}

func TestGetByUserIdLogs_Error(t *testing.T) {
	mockRepo := new(mocks.RentalLogsRepository)
	s := service.NewRentalLogsService(mockRepo)

	mockRepo.On("GetByUserId", 404).Return(nil, errors.New("no data"))

	resp, err := s.GetByUserIdLogs(404)

	assert.Error(t, err)
	assert.Len(t, resp, 0)
	assert.EqualError(t, err, "no data")
	mockRepo.AssertExpectations(t)
}
