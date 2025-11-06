package service

import (
	"log"
	"milestone2/internal/entity"
)

type CarsRepository interface {
	Create(car *entity.Cars) (err error)
	GetById(id int) (car entity.Cars, err error)
	GetAll() (cars []entity.Cars, err error)
}

type CarsServ struct {
	carsRepo CarsRepository
}

func NewCarsService(carsRepo CarsRepository) *CarsServ {
	return &CarsServ{carsRepo}
}

func (cs *CarsServ) Create(req entity.CreateRentalCarsRequest) (carResponse entity.CarsResponse, err error) {
	cars := entity.Cars{
		Name: req.Name,
		PlatNumber: req.PlatNumber,
		CategoryId: req.CategoryId,
		Description: req.Description,
		Price: req.Price,
		Availability: req.Availability,
	}
	if err := cs.carsRepo.Create(&cars); err != nil {
		log.Printf("error create cars on service %s", err)
		return entity.CarsResponse{}, err
	}	

	carsResponse, err := cs.GetById(cars.Id)
	if err != nil {
		log.Printf("error getting cars by id on service %s", err)
		return entity.CarsResponse{}, err
	}

	return carsResponse, nil
}

func (cs *CarsServ) GetById(id int) (carResponse entity.CarsResponse, err error) {
	car, err := cs.carsRepo.GetById(id)
	if err != nil {
		log.Printf("error getting cars by id on service %s", err)
		return entity.CarsResponse{}, err
	}

	carsResponse := entity.CarsResponse{
		Id: car.Id,
		Name: car.Name,
		PlatNumber: car.PlatNumber,
		CategoryId: car.CategoryId,
		Categories: car.Categories,
		Description: car.Description,
		Price: car.Price,
		Availability: car.Availability,
	}

	return carsResponse, nil
}

func (cs *CarsServ) GetAll() (carsResponse []entity.CarsResponse, err error) {
	cars, err := cs.carsRepo.GetAll()
	if err != nil {
		log.Printf("error get all on service %s", err)
		return []entity.CarsResponse{}, err
	}

	for _, car := range cars {
		carsResponse = append(carsResponse, entity.CarsResponse{
			Id: car.Id,
			Name: car.Name,
			PlatNumber: car.PlatNumber,
			CategoryId: car.CategoryId,
			Categories: car.Categories,
			Description: car.Description,
			Price: car.Price,
			Availability: car.Availability,
		})
	}
	
	return carsResponse, nil
}