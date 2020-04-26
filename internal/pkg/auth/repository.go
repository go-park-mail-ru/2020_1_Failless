package auth

import "failless/internal/pkg/models"

type Repository interface {
	GetUserByPhoneOrEmail(phone string, email string) (models.User, error)
}
