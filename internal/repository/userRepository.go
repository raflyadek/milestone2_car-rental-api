package repository

import (
	"context"
	"log"
	"milestone2/internal/entity"
	"os"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (ur *UserRepo) Create(user *entity.User) (err error) {
	if err := ur.db.WithContext(context.Background()).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) GetByEmail(email string) (user entity.User, err error) {
	if err := ur.db.WithContext(context.Background()).First(&user, "email = ?", email).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (ur *UserRepo) GetById(id int) (user entity.User, err error) {
	if err := ur.db.WithContext(context.Background()).First(&user, "id = ?", id).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (ur *UserRepo) UpdateValidationStatus(code, email string) (err error) {
	var user entity.User
	err = ur.db.WithContext(context.Background()).Model(&user).
	Where("validation_code = ? AND email = ?", code, email).
	Update("validation_status", true).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) SendValidationCode(send entity.SendEmailValidationRequest) (error) {
	keyPublic := os.Getenv("MAILJET_API_KEY")
	keyPrivate := os.Getenv("MAILJET_SECRET_KEY")
	mailjetClient := mailjet.NewMailjetClient(keyPublic, keyPrivate)
	messagesInfo := []mailjet.InfoMessagesV31 {
      mailjet.InfoMessagesV31{
        From: &mailjet.RecipientV31{
          Email: "carzrentalz@fivermail.com",
          Name: "Carz Rentalz",
        },
        To: &mailjet.RecipientsV31{
          mailjet.RecipientV31 {
            Email: send.Email,
            Name: send.Name,
          },
        },
        Subject: send.Subject,
        TextPart: send.TextPart,
        HTMLPart: send.HtmlPart,
      },
    }
	messages := mailjet.MessagesV31{Info: messagesInfo }
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}