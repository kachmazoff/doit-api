package impl

import (
	"fmt"

	"github.com/kachmazoff/doit-api/internal/mailing"
	"github.com/kachmazoff/doit-api/internal/model"
	"github.com/kachmazoff/doit-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repo   repository.Users
	sender mailing.Sender
}

func NewUsersService(repo repository.Users, sender mailing.Sender) *UsersService {
	return &UsersService{
		repo:   repo,
		sender: sender,
	}
}

func (u *UsersService) Create(newUser model.User) (string, error) {
	password := []byte(newUser.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 6)
	if err != nil {
		return "", err
	}
	newUser.Password = string(hashedPassword)

	userId, err := u.repo.Create(newUser)
	newUser.Id = userId

	if err == nil {
		go sendVerificationMessage(u.sender, newUser)
	}

	return userId, err
}

func (u *UsersService) ConfirmAccount(userId string) error {
	err := u.repo.SetStatus(userId, "active")
	if err == nil {
		var email string
		email, err = u.repo.GetEmailById(userId)
		if err == nil {
			go sendSuccessVerificationMessage(u.sender, email)
		}
	}

	return err
}

func (u *UsersService) GetIdByUsername(username string) (string, error) {
	return u.repo.GetIdByUsername(username)
}

func (u *UsersService) GetByUsername(username string) (model.User, error) {
	return u.repo.GetByUsername(username)
}

func (u *UsersService) GetByEmail(email string) (model.User, error) {
	return u.repo.GetByEmail(email)
}

func sendVerificationMessage(sender mailing.Sender, user model.User) {
	const template = "Вы зарегистрировались в сервисе <b>Doit</b>, как <i>%s</i>.<br></br> Для подтверждения регистрации перейдите по ссылке <a href='http://localhost:8080/api/auth/activate?id=%s'>Подтвердить</a>"
	mailBody := fmt.Sprintf(template, user.Username, user.Id)
	sender.Send(user.Email, "Регистрация в Doit", mailBody)
}

func sendSuccessVerificationMessage(sender mailing.Sender, email string) {
	const template = "Подтверждение регистрации прошло успешно!"
	sender.Send(email, "Регистрация в Doit", template)
}
