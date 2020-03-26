package usecase

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"log"
	"net/http"
)

func (uc *userUseCase) UpdateUserMeta(form *forms.MetaForm) (int, error) {
	if err := uc.rep.AddUserInfo(user, info); err != nil {
		return http.StatusNotModified, err
	}

	return 0, nil
}

func (uc *userUseCase) UpdateUserInfo(form *forms.GeneralForm) (int, error) {
	var info models.JsonInfo
	var user models.User

	if err := form.GetDBFormat(&info, &user); err != nil {
		return http.StatusInternalServerError, err
	}
	user.Uid = form.Uid
	if err := uc.rep.AddUserInfo(user, info); err != nil {
		return http.StatusNotModified, err
	}

	form.Avatar.ImgBase64 = ""
	for _, item := range form.Photos {
		item.ImgBase64 = ""
		item.Img = nil
	}
	return http.StatusOK, nil
}

func (uc *userUseCase) GetUserInfo(profile *forms.GeneralForm) (int, error) {
	row, err := uc.rep.GetProfileInfo(profile.Uid)
	if err != nil {
		log.Println(err.Error())
		return http.StatusNotFound, err
	}

	err = profile.FillProfile(row)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	base, err := uc.rep.GetUserByUID(profile.Uid)
	if err != nil {
		return http.StatusNotFound, err
	}

	profile.SignForm.Name = base.Name
	profile.Phone = base.Phone
	profile.Email = base.Email
	profile.Uid = base.Uid

	profile.Events, err = uc.rep.GetUserEvents(base.Uid)
	if err != nil {
		log.Println("error in get user events. Not fatal")
		log.Println(err.Error())
	}

	profile.Tags, err = uc.rep.GetUserTags(base.Uid)
	if err != nil {
		log.Println("error in get user tags. Not fatal")
		log.Println(err.Error())
	}

	return http.StatusOK, nil
}

func (uc *userUseCase) AddImageToProfile(uid int, name string) error {
	err := uc.rep.UpdateUserPhotos(uid, name)
	return err
}
