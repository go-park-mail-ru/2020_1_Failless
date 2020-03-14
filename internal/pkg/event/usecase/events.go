package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/repository"
	"failless/internal/pkg/models"
	"net/http"
)

type userUseCase struct {
	rep event.Repository
}

func GetUseCase() event.UseCase {
	return &userUseCase{
		rep: repository.NewSqlEventRepository(db.ConnectToDB()),
	}
}


func (uc *userUseCase) InitEventsByTime(events []models.Event) (int, error) {
	events, err := uc.rep.GetAllEvents()
	if err != nil {
		//network.GenErrorCode(w, r, err.Error(), )
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
