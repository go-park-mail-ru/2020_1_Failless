package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/repository"
)

type UseCase struct {
	rep user.Repository
}

func GetUseCase() UseCase {
	return UseCase{
		rep: repository.NewSqlUserRepository(db.ConnectToDB()),
	}
}
