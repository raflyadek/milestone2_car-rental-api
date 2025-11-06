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
	carsRepo := repository.NewCarsRepository(db)

	//service
	userServ := service.NewUserService(userRepo)
	carsServ := service.NewCarsService(carsRepo)

	//handler
	userHand := handler.NewUserHandler(userServ)
	carsHand := handler.NewCarsHandler(carsServ)

	//echo
	e := echo.New()
	e.Use(middleware.LoggingMiddleware)
	e.HTTPErrorHandler = middleware.ErrorHandler
	//auth
	e.POST("/users/register", userHand.UserRegister)
	e.POST("/users/login", userHand.UserLogin)
	e.PUT("/users/validation", userHand.UserValidation)
	
	//restricted endpoint 
	jwt := e.Group("")
	jwt.Use(middleware.LoggingMiddleware)
	jwt.Use(middleware.JwtMiddleware)
	//admin
	jwt.POST("/admin/cars", carsHand.CreateRentalCars)
	//all
	jwt.GET("/users/cars/:id", carsHand.GetRentalCarsById)
	jwt.GET("/users/cars", carsHand.GetAllCars)


	port := os.Getenv("PORT")
	if err := e.Start(":"+port); err != nil {
		logrus.Error("error connect to server", err.Error())
	}
}