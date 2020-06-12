package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/email"
	"failless/internal/pkg/email/repository"
)

type emailUseCase struct {
	rep email.Repository
}

func GetUseCase() email.UseCase {
	return &emailUseCase{
		rep: repository.NewSqlEmailRepository(db.ConnectToDB()),
	}
}
