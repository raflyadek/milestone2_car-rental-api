package handler

import (
	"milestone2/internal/entity"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type LogsService interface {
	CheckAvailabilityByCarId(req entity.CheckCarAvailabilityRequest) (resp entity.RentalAvailabilityResponse, err error)
}

type CarAvailLogsHandler struct {
	logsService LogsService
	validate *validator.Validate
}

func NewCarAvailLogsHandler(logsService LogsService, validate *validator.Validate) *CarAvailLogsHandler {
	return &CarAvailLogsHandler{logsService, validate}
}

func (al *CarAvailLogsHandler) CheckAvailabilityByCarId(c echo.Context) error {
	req := new(entity.CheckCarAvailabilityRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err := al.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	resp, err := al.logsService.CheckAvailabilityByCarId(*req)
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": resp,
	})
}