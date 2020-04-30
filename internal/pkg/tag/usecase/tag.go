package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/settings"
	"failless/internal/pkg/tag"
	"failless/internal/pkg/tag/repository"
	"log"
	"net/http"
)

type tagUseCase struct {
	rep tag.Repository
}

func GetUseCase() tag.UseCase {
	if settings.UseCaseConf.InHDD {
		return &tagUseCase{
			rep: repository.NewSqlTagRepository(db.ConnectToDB()),
		}
	} else {
		return &tagUseCase{
			rep: repository.NewTagRepository(),
		}
	}
}

func (uc *tagUseCase) InitEventsByTime(tags *models.TagList) (status int, err error) {
	*tags, err = uc.rep.GetAllTags()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	log.Println(tags)
	return http.StatusOK, nil
}
