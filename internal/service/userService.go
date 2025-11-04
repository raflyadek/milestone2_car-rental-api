package service

import (
	"log"
	"milestone2/internal/entity"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user entity.User) (register entity.UserRegister, err error)
	GetByEmail(email string) (user entity.User, err error)
}

type UserServ struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserServ {
	return &UserServ{userRepo}
}

func (us *UserServ) Create(user entity.User) (entity.UserRegister, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error generate hash on service %s", err)
		return entity.UserRegister{}, err
	}
	user.Password = string(hashPass)

	//create doesnt return anything but yeah just try i think 
	register, err := us.userRepo.Create(user); 
	if err != nil {
		log.Printf("error create user on service %s", err)
		return entity.UserRegister{}, err
	}

	return register, nil
}

func (us *UserServ) GetByEmail(email, password string) (accessToken string, err error) {
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