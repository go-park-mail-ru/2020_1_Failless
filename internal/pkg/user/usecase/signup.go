package usecase

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
)

func (uc *UseCase) RegisterNewUser(user *forms.SignForm) error {
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

	return uc.rep.AddNewUser(&dbUser)
}
