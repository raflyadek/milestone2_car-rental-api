package handler

import (
	"milestone2/internal/entity"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type RentalLogsService interface {
	GetAllLogs() (resp []entity.RentalLogsResponseAdmin, err error)
	GetByIdLogs(id int) (resp entity.RentalLogsResponseAdmin, err error)
	GetByUserIdLogs(userId int) (resp []entity.RentalLogsResponseUser, err error)
}

type RentalLogsHandler struct {
	logsServ RentalLogsService
}

func NewRentalLogsHandler(logsServ RentalLogsService) *RentalLogsHandler {
	return &RentalLogsHandler{logsServ}
}

func (rlh *RentalLogsHandler) GetAllLogs(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	role := claim["role"].(string)

	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "forbidden to access",
		})
	}

	logs, err := rlh.logsServ.GetAllLogs()
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data":    logs,
	})
}

func (rlh *RentalLogsHandler) GetByIdLogs(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	role := claim["role"].(string)

	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "forbidden to access",
		})
	}

	logs, err := rlh.logsServ.GetByIdLogs(id)
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data":    logs,
	})
}

func (rlh *RentalLogsHandler) GetByUserIdLogs(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	userId := int(claim["id"].(float64))

	logs, err := rlh.logsServ.GetByUserIdLogs(userId)
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": logs,
	})
}