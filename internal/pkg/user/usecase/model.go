package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/settings"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/repository"
	"log"
	"net/http"
)

type UserUseCase struct {
	Rep user.Repository
}

func GetUseCase() user.UseCase {
	if settings.UseCaseConf.InHDD {
		log.Println("IN HDD")
		return &UserUseCase{
			Rep: repository.NewSqlUserRepository(db.ConnectToDB()),
		}
	} else {
		log.Println("IN MEMORY")
		return &UserUseCase{
			Rep: repository.NewUserRepository(),
		}
	}
}

func (uc *UserUseCase) InitUsersByUserPreferences(users *[]models.UserGeneral, request *models.UserRequest) (int, error) {
	_, err := uc.Rep.GetValidTags()
	if err != nil {
		return http.StatusBadRequest, err
	}

	*users, err = uc.Rep.GetRandomFeedUsers(request.Uid, request.Limit, request.Page)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	log.Println(users)
	return http.StatusOK, nil
}

func (uc *UserUseCase) TakeValidTagsOnly(tagIds []int, tags []models.Tag) []int {
	var valid []int = nil
	for _, tagId := range tagIds {
		for _, tag := range tags {
			if tagId == tag.TagId {
				valid = append(valid, tagId)
			}
		}
	}

	return valid
}

func (uc *UserUseCase) GetUserSubscriptions(events *[]models.Event, uid int) (int, error) {
	var err error
	*events, err = uc.Rep.GetUserSubscriptions(uid)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
