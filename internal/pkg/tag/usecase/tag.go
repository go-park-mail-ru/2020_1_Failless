package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/tag"
	"failless/internal/pkg/tag/repository"
	"net/http"
)

type tagUseCase struct {
	rep tag.Repository
}

func GetUseCase() tag.UseCase {
	return &tagUseCase{
		rep: repository.NewSqlTagRepository(db.ConnectToDB()),
	}
}

func (uc *tagUseCase) InitEventsByTime(tags []models.Tag) (int, error) {
	tags, err := uc.rep.GetAllTags()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
