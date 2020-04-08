package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/settings"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/repository"
	"log"
)

type userUseCase struct {
	rep user.Repository
}

func GetUseCase() user.UseCase {
	if settings.UseCaseConf.InHDD {
		log.Println("IN HDD")
		return &userUseCase{
			rep: repository.NewSqlUserRepository(db.ConnectToDB()),
		}
	} else {
		log.Println("IN MEMORY")
		return &userUseCase{
			rep: repository.NewUserRepository(),
		}
	}
}
