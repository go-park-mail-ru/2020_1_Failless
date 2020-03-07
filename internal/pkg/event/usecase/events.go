package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/repository"
	"failless/internal/pkg/models"
	"net/http"
)

type UseCase struct {
	rep event.Repository
}

func GetUseCase() UseCase {
	return UseCase{
		rep: repository.NewSqlEventRepository(db.ConnectToDB()),
	}
}


func (uc *UseCase) InitEventsByTime(events []models.Event) (int, error) {
	events, err := uc.rep.GetAllEvents()
	if err != nil {
		//network.GenErrorCode(w, r, err.Error(), )
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
