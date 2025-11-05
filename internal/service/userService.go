package service

import (
	"crypto/rand"
	"fmt"
	"log"
	"milestone2/internal/entity"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *entity.User) (err error)
	GetByEmail(email string) (user entity.User, err error)
	GetById(id int) (user entity.User, err error)
	UpdateValidationStatus(code, email string) (err error)
}

type UserServ struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserServ {
	return &UserServ{userRepo}
}

func (us *UserServ) CreateUser(user entity.User) (userInfo entity.UserInfo, err error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error generate hash on service %s", err)
		return entity.UserInfo{}, err
	}
	user.Password = string(hashPass)

	//validate code
	secretCode := os.Getenv("CODE")
	codeInt, err := strconv.Atoi(secretCode)
	if err != nil {
		log.Printf("error casting string to int %s", err)
		return
	}

	code, err := us.generateValidationCode(codeInt)
	if err != nil {
		log.Printf("error generate validation code on service %s", err)
		return entity.UserInfo{}, err
	}
	user.ValidationCode = code

	//create doesnt return anything but yeah just try i think 
	if err := us.userRepo.Create(&user); err != nil {
		log.Printf("error create user on service %s", err)
		return entity.UserInfo{}, err
	}

	infoUser, err := us.GetUserById(user.Id)
	if err != nil {
		log.Printf("error get user info by id %s", err)
		return 
	}

	return infoUser, nil
}

func (us *UserServ) GetUserByEmail(email, password string) (accessToken string, err error) {
	user, err := us.userRepo.GetByEmail(email)
	if err != nil {
		log.Printf("failed get by email on service %s", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("failed compared hash password and passwor don service %s", err)
		return "", err
	}

	token, err := us.generateToken(user.Id, user.Email)
	if err != nil {
		log.Printf("error generate token on service %s", err)
		return "", err
	}

	return token, nil
}

func (us *UserServ) GetUserById(id int) (entity.UserInfo, error) {
	user, err := us.userRepo.GetById(id)
	if err != nil {
		log.Printf("failed get by id on service %s", err)
		return entity.UserInfo{}, err
	}

	userInfo := entity.UserInfo{
		Id: user.Id,
		Email: user.Email,
		FullName: user.FullName,
		Deposit: user.Deposit,
	}

	return userInfo, nil
}

func (us *UserServ) GetUserValidation(code, email string) (entity.UserInfo, error) {
	if err := us.userRepo.UpdateValidationStatus(code, email); err != nil {
		log.Printf("failed update validation status user on service %s", err)
		return entity.UserInfo{}, err
	}

	user, err := us.userRepo.GetByEmail(email)
	if err != nil {
		log.Printf("failed get user by email on service %s", err)
		return entity.UserInfo{}, err
	}

	userInfo := entity.UserInfo{
		Id: user.Id,
		Email: user.Email,
		FullName: user.FullName,
		Deposit: user.Deposit,
		ValidationStatus: user.ValidationStatus,
	}

	return userInfo, nil
}

func (us *UserServ) generateToken(id int, email string) (token string, err error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, 
	jwt.MapClaims{
		"id": id,
		"email": email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := jwt.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us *UserServ) generateValidationCode(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
