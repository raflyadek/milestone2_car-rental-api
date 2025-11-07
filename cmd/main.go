package main

import (
	"milestone2/config"
	"milestone2/internal/handler"
	"milestone2/internal/middleware"
	"milestone2/internal/repository"
	"milestone2/internal/service"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	db := config.ConnectionDb()
	validator := validator.New()

	//dependecy injection
	//repository
	userRepo := repository.NewUserRepository(db)
	carsRepo := repository.NewCarsRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	rentalLogsRepo := repository.NewRentalLogsRepository(db)

	//service
	userServ := service.NewUserService(userRepo)
	carsServ := service.NewCarsService(carsRepo)
	paymentServ := service.NewPaymentService(paymentRepo, carsRepo)
	rentalLogsServ := service.NewRentalLogsService(rentalLogsRepo)

	//handler
	userHand := handler.NewUserHandler(userServ)
	carsHand := handler.NewCarsHandler(carsServ)
	paymentHand := handler.NewPaymentHandler(paymentServ, validator)
	rentalLogsHand := handler.NewRentalLogsHandler(rentalLogsServ)

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
	jwt.PATCH("/admin/payments/:id", paymentHand.TransactionUpdatePayment)
	jwt.GET("/admin/payments", paymentHand.GetAllPayment)
	jwt.GET("/admin/payments/:id", paymentHand.GetByIdPayment)
	jwt.GET("/admin/rental/logs", rentalLogsHand.GetAllLogs)
	jwt.GET("/admin/rental/logs/:id", rentalLogsHand.GetByIdLogs)
	//all
	jwt.GET("/users/cars/:id", carsHand.GetRentalCarsById)
	jwt.GET("/users/cars", carsHand.GetAllCars)
	jwt.POST("/users/payments", paymentHand.CreatePayment)
	jwt.GET("/users/payments", paymentHand.GetByUserIdPayment)
	jwt.GET("/users/rental/logs", rentalLogsHand.GetByUserIdLogs)


	port := os.Getenv("PORT")
	if err := e.Start(":"+port); err != nil {
		logrus.Error("error connect to server", err.Error())
	}
}