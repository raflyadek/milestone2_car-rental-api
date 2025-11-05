package main

import (
	"milestone2/config"
	"milestone2/internal/handler"
	"milestone2/internal/middleware"
	"milestone2/internal/repository"
	"milestone2/internal/service"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	db := config.ConnectionDb()

	//dependecy injection
	//repository
	userRepo := repository.NewUserRepository(db)

	//service
	userServ := service.NewUserService(userRepo)

	//handler
	userHand := handler.NewUserHandler(userServ)

	//echo
	e := echo.New()
	e.Use(middleware.LoggingMiddleware)
	e.HTTPErrorHandler = middleware.ErrorHandler
	//auth
	e.POST("/users/register", userHand.UserRegister)
	e.POST("/users/login", userHand.UserLogin)

	jwt := e.Group("")
	jwt.Use(middleware.LoggingMiddleware)
	jwt.Use(middleware.JwtMiddleware)


	port := os.Getenv("PORT")
	if err := e.Start(":"+port); err != nil {
		logrus.Error("error connect to server", err.Error())
	}
}