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

func (uc *UserUseCase) InitEventsByUserPreferences(users *[]models.UserGeneral, request *models.UserRequest) (int, error) {
	_, err := uc.Rep.GetValidTags()
	if err != nil {
		return http.StatusBadRequest, err
	}

	//valid := uc.TakeValidTagsOnly(request.Tags, dbTags)
	//log.Println(request)
	//if valid != nil {gst
	//	*users, err = uc.Rep.GetNewUsersByTags(valid, request.Uid, request.Limit, request.Page)
	//} else {
	//	*users, err = uc.Rep.GetRandomFeedUsers(request.Uid, request.Limit, request.Page)
	//}
	*users, err = uc.Rep.GetRandomFeedUsers(request.Uid, request.Limit, request.Page)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	//for i := 0; i < len(*users); i++ {
	//	(*users)[i].Tag = dbTags[(*users)[i].Type-1]
	//	log.Println(dbTags[(*users)[i].Type - 1])
	//}

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
