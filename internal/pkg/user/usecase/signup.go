package usecase

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"log"
)

func (uc *UserUseCase) RegisterNewUser(user *forms.SignForm) error {
	// TODO: move it to repository
	bPass, err := security.EncryptPassword(user.Password)
	if err != nil {
		return err
	}

	dbUser := models.User{
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: bPass,
	}
	log.Println("ORIGINAL: ", user.Password)
	log.Println("ORIGINAL: ", bPass)
	return uc.Rep.AddNewUser(&dbUser)
}
