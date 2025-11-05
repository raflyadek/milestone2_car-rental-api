package handler

import (
	"milestone2/internal/entity"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	CreateUser(user entity.User) (userInfo entity.UserInfo, err error) 
	GetUserByEmail(email, password string) (accessToken string, err error)
	GetUserById(id int) (entity.UserInfo, error)
	GetUserValidation(code, email string) (entity.UserInfo, error)
}

type UserHandler struct {
	userServ UserService
}

func NewUserHandler(userServ UserService) *UserHandler {
	return &UserHandler{userServ}
}

func (uh *UserHandler) UserRegister(c echo.Context) error {
	req := new(entity.RegisterRequest)
	if err := c.Bind(&req); err != nil {
		logrus.Error("failed bind register request on handler", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		logrus.Error("failed validate register request on handler", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": err.Error(),
		})
	}

	userInfo, err := uh.userServ.CreateUser(entity.User{
		Email: req.Email,
		FullName: req.FullName,
		Password: req.Password,
	})
	if err != nil {
		logrus.Error("failed execute create user on handler", err.Error())
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "register failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success register user, check your email for email validation",
		"data": userInfo,
	})
}

func (uh *UserHandler) UserLogin(c echo.Context) error {
	req := new(entity.LoginRequest)
	if err := c.Bind(&req); err != nil {
		logrus.Error("failed bind login request on handler", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		logrus.Error("failed validate login request on handler", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": err.Error(),
		})
	}	

	jwtToken, err := uh.userServ.GetUserByEmail(req.Email, req.Password)
	if err != nil {
		logrus.Error("failed execute get by email on handler", err.Error())
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "login failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success login user",
		"data": jwtToken,
	})
}

func (uh *UserHandler) UserValidation(c echo.Context) error {
	req := new(entity.EmailValidationRequest)	
	if err := c.Bind(&req); err != nil {
		logrus.Error("failed bind login request on handler", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		logrus.Error("failed validate login request on handler", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": err.Error(),
		})
	}

	user, err := uh.userServ.GetUserValidation(req.Code, req.Email)
	if err != nil {
		logrus.Error("failed execute get by validation code on handler", err.Error())
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"message": "failed validate the email",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success validate user",
		"data": user,
	})
}