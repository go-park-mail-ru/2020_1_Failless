package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/repository"
)

type userUseCase struct {
	rep user.Repository
}

func GetUseCase() user.UseCase {
	return &userUseCase{
		rep: repository.NewSqlUserRepository(db.ConnectToDB()),
	}
}
