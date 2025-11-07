package handler_test

import (
	"errors"
	"milestone2/internal/entity"
	"milestone2/internal/handler"
	"milestone2/internal/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupEchoContext(tokenClaims jwt.MapClaims) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// inject mock user token
	token := &jwt.Token{Claims: tokenClaims}
	c.Set("user", token)

	return c, rec
}

func TestGetAllLogs_AdminAccess(t *testing.T) {
	mockServ := new(mocks.RentalLogsService)
	handler := handler.NewRentalLogsHandler(mockServ)

	// expected data
	expectedLogs := []entity.RentalLogsResponseAdmin{{Id: 1, UserId: 2}}
	mockServ.On("GetAllLogs").Return(expectedLogs, nil)

	c, rec := setupEchoContext(jwt.MapClaims{"role": "admin"})

	err := handler.GetAllLogs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "ok")
	assert.Contains(t, rec.Body.String(), `"id":1`)
	mockServ.AssertExpectations(t)
}

func TestGetAllLogs_ForbiddenForUser(t *testing.T) {
	mockServ := new(mocks.RentalLogsService)
	handler := handler.NewRentalLogsHandler(mockServ)

	c, rec := setupEchoContext(jwt.MapClaims{"role": "user"})

	err := handler.GetAllLogs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "forbidden")
}

func TestGetByIdLogs_AdminOK(t *testing.T) {
	mockServ := new(mocks.RentalLogsService)
	handler := handler.NewRentalLogsHandler(mockServ)

	expected := entity.RentalLogsResponseAdmin{Id: 1, UserId: 2}
	mockServ.On("GetByIdLogs", 1).Return(expected, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/logs/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", &jwt.Token{Claims: jwt.MapClaims{"role": "admin"}})

	err := handler.GetByIdLogs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":1`)
	mockServ.AssertExpectations(t)
}

func TestGetByUserIdLogs_UserOK(t *testing.T) {
	mockServ := new(mocks.RentalLogsService)
	handler := handler.NewRentalLogsHandler(mockServ)

	expectedLogs := []entity.RentalLogsResponseUser{{Id: 99, CarId: 100}}
	mockServ.On("GetByUserIdLogs", 42).Return(expectedLogs, nil)

	c, rec := setupEchoContext(jwt.MapClaims{"id": float64(42), "role": "user"})

	err := handler.GetByUserIdLogs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":99`)
	mockServ.AssertExpectations(t)
}

func TestGetByUserIdLogs_Error(t *testing.T) {
	mockServ := new(mocks.RentalLogsService)
	handler := handler.NewRentalLogsHandler(mockServ)

	mockServ.On("GetByUserIdLogs", 5).Return(nil, errors.New("db error"))

	c, rec := setupEchoContext(jwt.MapClaims{"id": float64(5), "role": "user"})

	err := handler.GetByUserIdLogs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "internal server error")
}
