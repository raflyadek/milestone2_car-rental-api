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

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupEchoRequest(method, path string, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var req *http.Request

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func createToken(role string, id int) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"id":   float64(id),
	})
}

func TestCreatePayment_Success(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodPost, "/payments", entity.CreatePaymentRequest{
		CarId: 1, StartDate: "2025-11-01", EndDate: "2025-11-05",
	})
	token := createToken("user", 101)
	c.Set("user", token)

	mockServ := new(mocks.PaymentService)
	validate := validator.New()

	expectedResp := entity.PaymentInfoResponse{
		Id: 10, UserId: 101, CarId: 1, Price: 500000, Status: true,
	}
	mockServ.On("CreatePayment", 101, mock.AnythingOfType("entity.CreatePaymentRequest")).
		Return(expectedResp, nil)

	handler := NewPaymentHandler(mockServ, validate)

	err := handler.CreatePayment(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"ok"`)
	assert.Contains(t, rec.Body.String(), `"payment_id":10`)
	mockServ.AssertExpectations(t)
}

func TestCreatePayment_BadRequest(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodPost, "/payments", "invalid-json")
	token := createToken("user", 1)
	c.Set("user", token)

	mockServ := new(mocks.PaymentService)
	validate := validator.New()
	handler := NewPaymentHandler(mockServ, validate)

	err := handler.CreatePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAllPayment_AdminOnly(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodGet, "/payments", nil)
	token := createToken("user", 1) // not admin
	c.Set("user", token)

	mockServ := new(mocks.PaymentService)
	validate := validator.New()
	handler := NewPaymentHandler(mockServ, validate)

	err := handler.GetAllPayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGetAllPayment_Success(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodGet, "/payments", nil)
	token := createToken("admin", 1)
	c.Set("user", token)

	mockServ := new(mocks.PaymentService)
	validate := validator.New()

	mockServ.On("GetAllPayment").Return([]entity.PaymentInfoResponse{
		{Id: 1, UserId: 1, CarId: 10},
	}, nil)

	handler := NewPaymentHandler(mockServ, validate)
	err := handler.GetAllPayment(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"payment_id":1`)
	mockServ.AssertExpectations(t)
}

func TestGetByUserIdPayment_Success(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodGet, "/payments/user", nil)
	token := createToken("user", 99)
	c.Set("user", token)

	mockServ := new(mocks.PaymentService)
	validate := validator.New()

	mockServ.On("GetByUserIdPayment", 99).Return([]entity.PaymentInfoResponse{
		{Id: 5, UserId: 99, CarId: 7},
	}, nil)

	handler := NewPaymentHandler(mockServ, validate)
	err := handler.GetByUserIdPayment(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"payment_id":5`)
	mockServ.AssertExpectations(t)
}

func TestGetByIdPayment_Success(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodGet, "/payments/3", nil)
	c.SetParamNames("id")
	c.SetParamValues("3")

	mockServ := new(mocks.PaymentService)
	validate := validator.New()

	mockServ.On("GetByIdPayment", 3).Return(entity.PaymentInfoResponse{
		Id: 3, UserId: 1, CarId: 5,
	}, nil)

	handler := NewPaymentHandler(mockServ, validate)
	err := handler.GetByIdPayment(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"payment_id":3`)
	mockServ.AssertExpectations(t)
}

func TestTransactionUpdatePayment_AdminOnly(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodPut, "/payments/1", nil)
	token := createToken("user", 1)
	c.Set("user", token)

	mockServ := new(mocks.PaymentService)
	validate := validator.New()
	handler := NewPaymentHandler(mockServ, validate)

	err := handler.TransactionUpdatePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestTransactionUpdatePayment_Success(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodPut, "/payments/1", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")
	token := createToken("admin", 1)
	c.Set("user", token)

	mockServ := new(mocks.PaymentService)
	validate := validator.New()

	mockServ.On("TransactionUpdatePayment", 1).Return(entity.PaidPaymentResponse{
		Id: 1, UserId: 1, CarId: 10, TotalSpent: 700000,
	}, nil)

	handler := NewPaymentHandler(mockServ, validate)
	err := handler.TransactionUpdatePayment(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":1`)
	mockServ.AssertExpectations(t)
}

func TestTransactionUpdatePayment_Error(t *testing.T) {
	c, rec := setupEchoRequest(http.MethodPut, "/payments/1", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")
	token := createToken("admin", 1)
	c.Set("user", token)

	mockServ := new(mocks.PaymentService)
	validate := validator.New()

	mockServ.On("TransactionUpdatePayment", 1).Return(entity.PaidPaymentResponse{}, errors.New("fail"))

	handler := NewPaymentHandler(mockServ, validate)
	err := handler.TransactionUpdatePayment(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockServ.AssertExpectations(t)
}
