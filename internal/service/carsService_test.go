package service_test

import (
	"errors"
	"milestone2/internal/entity"
	"milestone2/internal/service"
	"milestone2/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCarsServ_Create(t *testing.T) {
	mockRepo := new(mocks.CarsRepository)
	carsServ := service.NewCarsService(mockRepo)

	req := entity.CreateRentalCarsRequest{
		Name:         "Toyota Avanza",
		PlatNumber:   "B1234XYZ",
		CategoryId:   1,
		Description:  "Family car",
		Price:        300000,
		Availability: true,
	}

	car := entity.Cars{
		Id:           1,
		Name:         req.Name,
		PlatNumber:   req.PlatNumber,
		CategoryId:   req.CategoryId,
		Description:  req.Description,
		Price:        req.Price,
		Availability: req.Availability,
	}

	carResponse := entity.CarsResponse{
		Id:           car.Id,
		Name:         car.Name,
		PlatNumber:   car.PlatNumber,
		CategoryId:   car.CategoryId,
		Description:  car.Description,
		Price:        car.Price,
		Availability: car.Availability,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*entity.Cars")).Return(nil).Once()
		mockRepo.On("GetById", mock.AnythingOfType("int")).Return(car, nil).Once()

		res, err := carsServ.Create(req)

		assert.NoError(t, err)
		assert.Equal(t, carResponse, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed create", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*entity.Cars")).Return(errors.New("db error")).Once()

		res, err := carsServ.Create(req)

		assert.Error(t, err)
		assert.Equal(t, entity.CarsResponse{}, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get by id", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*entity.Cars")).Return(nil).Once()
		mockRepo.On("GetById", mock.AnythingOfType("int")).Return(entity.Cars{}, errors.New("not found")).Once()

		res, err := carsServ.Create(req)

		assert.Error(t, err)
		assert.Equal(t, entity.CarsResponse{}, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarsServ_GetById(t *testing.T) {
	mockRepo := new(mocks.CarsRepository)
	carsServ := service.NewCarsService(mockRepo)

	car := entity.Cars{
		Id:           1,
		Name:         "Civic",
		PlatNumber:   "B9999ZZZ",
		CategoryId:   2,
		Description:  "Sport car",
		Price:        500000,
		Availability: true,
	}
	carResponse := entity.CarsResponse{
		Id:           car.Id,
		Name:         car.Name,
		PlatNumber:   car.PlatNumber,
		CategoryId:   car.CategoryId,
		Description:  car.Description,
		Price:        car.Price,
		Availability: car.Availability,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetById", 1).Return(car, nil).Once()

		res, err := carsServ.GetById(1)

		assert.NoError(t, err)
		assert.Equal(t, carResponse, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		mockRepo.On("GetById", 2).Return(entity.Cars{}, errors.New("not found")).Once()

		res, err := carsServ.GetById(2)

		assert.Error(t, err)
		assert.Equal(t, entity.CarsResponse{}, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarsServ_GetAll(t *testing.T) {
	mockRepo := new(mocks.CarsRepository)
	carsServ := service.NewCarsService(mockRepo)

	cars := []entity.Cars{
		{
			Id:           1,
			Name:         "Avanza",
			PlatNumber:   "B1234ABC",
			CategoryId:   1,
			Description:  "Family car",
			Price:        300000,
			Availability: true,
		},
	}

	carsResponse := []entity.CarsResponse{
		{
			Id:           1,
			Name:         "Avanza",
			PlatNumber:   "B1234ABC",
			CategoryId:   1,
			Description:  "Family car",
			Price:        300000,
			Availability: true,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAll").Return(cars, nil).Once()

		res, err := carsServ.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, carsResponse, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		mockRepo.On("GetAll").Return([]entity.Cars{}, errors.New("db error")).Once()

		res, err := carsServ.GetAll()

		assert.Error(t, err)
		assert.Empty(t, res)
		mockRepo.AssertExpectations(t)
	})
}
