package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/tag"
	"failless/internal/pkg/tag/repository"
	"net/http"
)

type UseCase struct {
	rep tag.Repository
}

func GetUseCase() UseCase {
	return UseCase{
		rep: repository.NewSqlTagRepository(db.ConnectToDB()),
	}
}

func (uc *UseCase) InitEventsByTime(tags []models.Tag) (int, error) {
	tags, err := uc.rep.GetAllTags()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
