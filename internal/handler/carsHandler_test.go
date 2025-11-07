package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"milestone2/internal/entity"
	"milestone2/internal/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateRentalCars_Success(t *testing.T) {
	e := echo.New()
	reqBody := entity.CreateRentalCarsRequest{
		Name:       "Avanza",
		PlatNumber: "B1234CD",
		CategoryId: 1,
		Description: "Family car",
		Price:      200000,
		Availability: true,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/cars", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"role": "admin"})
	c.Set("user", token)

	mockServ := new(mocks.CarsService)
	mockServ.On("Create", reqBody).Return(entity.CarsResponse{
		Id: 1, Name: "Avanza", PlatNumber: "B1234CD",
		CategoryId: 1, Description: "Family car", Price: 200000, Availability: true,
	}, nil)

	handler := NewCarsHandler(mockServ)

	err := handler.CreateRentalCars(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "success create car for rental")

	mockServ.AssertExpectations(t)
}

func TestGetAllCars_Error(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cars", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockServ := new(mocks.CarsService)
	mockServ.On("GetAll").Return(nil, errors.New("db fail"))

	handler := NewCarsHandler(mockServ)
	err := handler.GetAllCars(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "internal server error")

	mockServ.AssertExpectations(t)
}
