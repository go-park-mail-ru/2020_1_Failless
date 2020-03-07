package usecase

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"net/http"
)

func (uc *UseCase) UpdateUserInfo(form *forms.ProfileForm) (error, int) {
	var info models.JsonInfo
	var user models.User

	if err := form.GetDBFormat(&info, &user); err != nil {
		return err, http.StatusInternalServerError
	}
	user.Uid = form.Uid
	if err := uc.rep.AddUserInfo(user, info); err != nil {
		return err, http.StatusNotModified
	}

	form.Avatar.ImgBase64 = ""
	for _, item := range form.Photos {
		item.ImgBase64 = ""
		item.Img = nil
	}
	return nil, http.StatusOK
}