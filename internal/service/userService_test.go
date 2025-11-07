package service_test

import (
	"errors"
	"milestone2/internal/entity"
	"milestone2/internal/mocks"
	"milestone2/internal/service"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	os.Setenv("CODE", "12345")

	user := entity.User{
		Id:       1,
		Email:    "test@mail.com",
		FullName: "Tester",
		Password: "plainpassword",
	}

	mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)
	mockRepo.On("GetById", 1).Return(entity.User{
		Id:       1,
		Email:    "test@mail.com",
		FullName: "Tester",
	}, nil)
	mockRepo.On("SendValidationCode", mock.Anything).Return(nil)

	resp, err := userServ.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, "test@mail.com", resp.Email)
	assert.Equal(t, "Tester", resp.FullName)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_CreateError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)
	os.Setenv("CODE", "12345")

	user := entity.User{
		Id:       1,
		Email:    "fail@mail.com",
		FullName: "ErrorUser",
		Password: "plainpassword",
	}

	mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(errors.New("db error"))

	resp, err := userServ.CreateUser(user)
	assert.Error(t, err)
	assert.Equal(t, entity.UserResponse{}, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	password := "secret"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := entity.User{
		Id:               1,
		Email:            "user@mail.com",
		Password:         string(hashed),
		ValidationStatus: true,
		Role:             "user",
	}

	mockRepo.On("GetByEmail", "user@mail.com").Return(user, nil)

	token, err := userServ.GetUserByEmail("user@mail.com", password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_NotValidated(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	user := entity.User{
		Id:               1,
		Email:            "no@mail.com",
		Password:         "whatever",
		ValidationStatus: false,
	}

	mockRepo.On("GetByEmail", "no@mail.com").Return(user, nil)

	token, err := userServ.GetUserByEmail("no@mail.com", "secret")

	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := entity.User{
		Id:               1,
		Email:            "user@mail.com",
		Password:         string(hashed),
		ValidationStatus: true,
	}

	mockRepo.On("GetByEmail", "user@mail.com").Return(user, nil)

	token, err := userServ.GetUserByEmail("user@mail.com", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	mockRepo.On("GetById", 1).Return(entity.User{
		Id:       1,
		Email:    "id@mail.com",
		FullName: "IDUser",
	}, nil)

	userResp, err := userServ.GetUserById(1)

	assert.NoError(t, err)
	assert.Equal(t, "id@mail.com", userResp.Email)
	assert.Equal(t, "IDUser", userResp.FullName)
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_Error(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	mockRepo.On("GetById", 1).Return(entity.User{}, errors.New("not found"))

	resp, err := userServ.GetUserById(1)

	assert.Error(t, err)
	assert.Equal(t, entity.UserResponse{}, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetUserValidation_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	user := entity.User{
		Id:               1,
		Email:            "user@mail.com",
		FullName:         "User",
		ValidationStatus: false,
	}

	mockRepo.On("GetByEmail", "user@mail.com").Return(user, nil)
	mockRepo.On("UpdateValidationStatus", "code123", "user@mail.com").Return(nil)

	resp, err := userServ.GetUserValidation("code123", "user@mail.com")

	assert.NoError(t, err)
	assert.Equal(t, true, resp.ValidationStatus)
	mockRepo.AssertExpectations(t)
}

func TestGetUserValidation_AlreadyValidated(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	user := entity.User{
		Id:               1,
		Email:            "user@mail.com",
		ValidationStatus: true,
	}

	mockRepo.On("GetByEmail", "user@mail.com").Return(user, nil)

	resp, err := userServ.GetUserValidation("code123", "user@mail.com")

	assert.Error(t, err)
	assert.Equal(t, entity.UserResponse{}, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetUserValidation_UpdateError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userServ := service.NewUserService(mockRepo)

	user := entity.User{
		Id:               1,
		Email:            "user@mail.com",
		ValidationStatus: false,
	}

	mockRepo.On("GetByEmail", "user@mail.com").Return(user, nil)
	mockRepo.On("UpdateValidationStatus", "code123", "user@mail.com").Return(errors.New("update failed"))

	resp, err := userServ.GetUserValidation("code123", "user@mail.com")

	assert.Error(t, err)
	assert.Equal(t, entity.UserResponse{}, resp)
	mockRepo.AssertExpectations(t)
}
