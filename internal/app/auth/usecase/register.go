package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/security"
)

func RegisterNewUser(user forms.SignForm) error {
	bPass, err := security.EncryptPassword(user.Password)
	if err != nil {
		return err
	}

	dbUser := db.User{
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: bPass,
	}

	return db.AddNewUser(db.ConnectToDB(), &dbUser)
}
