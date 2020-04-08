package usecase

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"log"
	"net/http"
)

func (uc *userUseCase) UpdateUserMeta(form *forms.MetaForm) (int, error) {
	if err := uc.rep.UpdateUserSimple(form.Uid, form.Social, &form.About); err != nil {
		return http.StatusNotModified, err
	}

	for _, tag := range form.Tags {
		if err := uc.rep.UpdateUserTags(form.Uid, tag); err != nil {
			return http.StatusNotModified, err
		}
	}
	return http.StatusOK, nil
}

func (uc *userUseCase) UpdateUserInfo(form *forms.GeneralForm) (int, error) {
	form.Photos = append(form.Photos, form.Avatar)

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

	(*profile).SignForm.Name = base.Name
	(*profile).Phone = base.Phone
	(*profile).Email = base.Email
	(*profile).Uid = base.Uid

	(*profile).Events, err = uc.rep.GetUserEvents(base.Uid)
	if err != nil {
		log.Println("error in get user events. Not fatal")
		log.Println(err.Error())
	}

	(*profile).Tags, err = uc.rep.GetUserTags(base.Uid)
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

func (uc *userUseCase) UpdateUserBase(form *forms.SignForm) (int, error) {
	usr, err := uc.rep.GetUserByUID(form.Uid)
	if err != nil {
		log.Println(err.Error())
		return http.StatusNotFound, err
	}

	usr.Uid = form.Uid
	//usr.Name = form.Name
	usr.Email = form.Email
	usr.Phone = form.Phone
	usr.Password, err = security.EncryptPassword(form.Password)

	var inf = models.JsonInfo{}
	//inf.Birthday = form.Birthday
	//inf.Gender = form.Gender

	if err := uc.rep.UpdUserGeneral(inf, usr); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
