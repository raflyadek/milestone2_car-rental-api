package handler_test

import (
	"errors"
	"milestone2/internal/entity"
	"milestone2/internal/handler"
	"milestone2/internal/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRegister_Success(t *testing.T) {
	mockServ := new(mocks.UserService)
	h := handler.NewUserHandler(mockServ)

	reqBody := `{"email":"test@mail.com","full_name":"Test User","password":"12345"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	expectedResp := entity.UserResponse{Id: 1, Email: "test@mail.com", FullName: "Test User"}
	mockServ.On("CreateUser", mock.AnythingOfType("entity.User")).Return(expectedResp, nil)

	err := h.UserRegister(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success register user`)
	mockServ.AssertExpectations(t)
}

func TestUserRegister_BindError(t *testing.T) {
	mockServ := new(mocks.UserService)
	h := handler.NewUserHandler(mockServ)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{invalid-json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.UserRegister(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "message")
}

func TestUserRegister_ServiceError(t *testing.T) {
	mockServ := new(mocks.UserService)
	h := handler.NewUserHandler(mockServ)

	reqBody := `{"email":"fail@mail.com","full_name":"Fail","password":"bad"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockServ.On("CreateUser", mock.AnythingOfType("entity.User")).Return(entity.UserResponse{}, errors.New("db error"))

	err := h.UserRegister(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "register failed")
}

func TestUserLogin_Success(t *testing.T) {
	mockServ := new(mocks.UserService)
	h := handler.NewUserHandler(mockServ)

	reqBody := `{"email":"test@mail.com","password":"12345"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockServ.On("GetUserByEmail", "test@mail.com", "12345").Return("mock-token", nil)

	err := h.UserLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "mock-token")
	mockServ.AssertExpectations(t)
}

func TestUserLogin_Error(t *testing.T) {
	mockServ := new(mocks.UserService)
	h := handler.NewUserHandler(mockServ)

	reqBody := `{"email":"wrong@mail.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockServ.On("GetUserByEmail", "wrong@mail.com", "wrong").Return("", errors.New("not found"))

	err := h.UserLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "login failed")
}

func TestUserValidation_Success(t *testing.T) {
	mockServ := new(mocks.UserService)
	h := handler.NewUserHandler(mockServ)

	reqBody := `{"code":"abc123","email":"valid@mail.com"}`
	req := httptest.NewRequest(http.MethodPost, "/validation", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	expectedUser := entity.UserResponse{Id: 1, Email: "valid@mail.com", FullName: "Valid User", ValidationStatus: true}
	mockServ.On("GetUserValidation", "abc123", "valid@mail.com").Return(expectedUser, nil)

	err := h.UserValidation(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success validate user`)
	mockServ.AssertExpectations(t)
}

func TestUserValidation_Error(t *testing.T) {
	mockServ := new(mocks.UserService)
	h := handler.NewUserHandler(mockServ)

	reqBody := `{"code":"xyz","email":"fail@mail.com"}`
	req := httptest.NewRequest(http.MethodPost, "/validation", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockServ.On("GetUserValidation", "xyz", "fail@mail.com").Return(entity.UserResponse{}, errors.New("invalid code"))

	err := h.UserValidation(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "validate failed")
}
