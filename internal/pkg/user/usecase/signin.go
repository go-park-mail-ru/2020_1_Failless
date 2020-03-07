package usecase

import (
	"errors"
	"failless/internal/pkg/models"
	"log"
	"net/http"
)

func (uc *UseCase) FillFormIfExist(cred *models.User) (int, error) {
	log.Println(*cred)
	user, err := uc.rep.GetUserByPhoneOrEmail(cred.Phone, cred.Email)
	if err == nil && user.Uid < 0 {
		log.Println("user not found")
		return http.StatusNotFound, errors.New("User doesn't exist\n")
	} else if err != nil {
		log.Println("error was occurred")
		log.Println(err.Error())
		return http.StatusInternalServerError, err
	}

	*cred = user
	return http.StatusOK, nil
}
