package handler

import (
	"log"
	"milestone2/internal/entity"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CarsService interface {
	Create(req entity.CreateRentalCarsRequest) (carResponse entity.CarsResponse, err error)
	GetById(id int) (carInfo entity.CarsResponse, err error)
	GetAll() (carsResponse []entity.CarsResponse, err error)
}

type CarsHandler struct {
	carsServ CarsService
}

func NewCarsHandler(carsServ CarsService) *CarsHandler {
	return &CarsHandler{carsServ}
}

func (ch *CarsHandler) CreateRentalCars(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	role := claim["role"].(string)

	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "you have no access ",
		})
	}

	req := new(entity.CreateRentalCarsRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	carInfo, err := ch.carsServ.Create(*req)
	if err != nil {
		log.Print(err.Error())
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create car for rental",
		"data": carInfo,
	})
}

func (ch *CarsHandler) GetRentalCarsById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	carInfo, err := ch.carsServ.GetById(id)
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": carInfo,
	})
}

func (ch *CarsHandler) GetAllCars(c echo.Context) error {
	cars, err := ch.carsServ.GetAll()
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": cars,
	})
}