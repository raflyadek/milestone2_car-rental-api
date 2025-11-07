package handler

import (
	"milestone2/internal/entity"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PaymentService interface {
	CreatePayment(userId int, req entity.CreatePaymentRequest) (resp entity.PaymentInfoResponse, err error)
	GetAllPayment() (resp []entity.PaymentInfoResponse, err error)
	GetByUserIdPayment(userId int) (resp []entity.PaymentInfoResponse, err error)
	GetByIdPayment(id int) (resp entity.PaymentInfoResponse, err error)
	TransactionUpdatePayment(paymentId int) (resp entity.PaidPaymentResponse, err error)
}

type PaymentHandler struct {
	paymentServ PaymentService
	validate *validator.Validate
}

func NewPaymentHandler(paymentServ PaymentService, validate *validator.Validate) *PaymentHandler {
	return &PaymentHandler{paymentServ, validate}
}

func (ph *PaymentHandler) CreatePayment(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	userId := int(claim["id"].(float64))

	req := new(entity.CreatePaymentRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err := ph.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	resp, err := ph.paymentServ.CreatePayment(userId, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": resp,
	})
}

func (ph *PaymentHandler) GetAllPayment(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	role := claim["role"].(string)

	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "forbidden to access",
		})
	}

	payments, err := ph.paymentServ.GetAllPayment()
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": payments,
	})
}

func (ph *PaymentHandler) GetByUserIdPayment(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	userId := int(claim["id"].(float64))

	payments, err := ph.paymentServ.GetByUserIdPayment(userId)
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": payments,
	})
}

func (ph *PaymentHandler) GetByIdPayment(c echo.Context) error { 
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	resp, err := ph.paymentServ.GetByIdPayment(id)
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

//it should be automate 
func (ph *PaymentHandler) TransactionUpdatePayment(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	role := claim["role"].(string)

	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "forbidden to access",
		})
	}

	idStr := c.Param("id")
	paymentId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	resp, err := ph.paymentServ.TransactionUpdatePayment(paymentId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		})		
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": resp,
	})
}