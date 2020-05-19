package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/tag"
	"failless/internal/pkg/tag/repository"
	"net/http"
)

type tagUseCase struct {
	Rep tag.Repository
}

func GetUseCase() tag.UseCase {
	return &tagUseCase{
		Rep: repository.NewSqlTagRepository(db.ConnectToDB()),
	}
}

func (uc *tagUseCase) InitEventsByTime(tags *models.TagList) (status int, err error) {
	*tags, err = uc.Rep.GetAllTags()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
