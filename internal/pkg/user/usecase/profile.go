package usecase

import (
	"failless/internal/pkg/db"
	eventRep "failless/internal/pkg/event/repository"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"log"
	"net/http"
)

func (uc *UserUseCase) UpdateUserAbout(uid int, about string) models.WorkMessage {
	return uc.Rep.UpdateUserAbout(uid, about)
}

func (uc *UserUseCase) UpdateUserTags(uid int, tagIDs []int) models.WorkMessage {
	return uc.Rep.UpdateUserTags(uid, tagIDs)
}

func (uc *UserUseCase) UpdateUserInfo(form *forms.GeneralForm) (int, error) {
	form.Photos = append(form.Photos, form.Avatar)

	var info models.JsonInfo
	var user models.User

	if err := form.GetDBFormat(&info, &user); err != nil {
		return http.StatusInternalServerError, err
	}
	user.Uid = form.Uid
	if err := uc.Rep.AddUserInfo(user, info); err != nil {
		return http.StatusNotModified, err
	}

	form.Avatar.ImgBase64 = ""
	for _, item := range form.Photos {
		item.ImgBase64 = ""
		item.Img = nil
	}
	return http.StatusOK, nil
}

func (uc *UserUseCase) GetUserInfo(profile *forms.GeneralForm) (int, error) {
	row, err := uc.Rep.GetProfileInfo(profile.Uid)
	if err != nil {
		log.Println(err.Error())
		return http.StatusNotFound, err
	}

	err = profile.FillProfile(row)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	base, err := uc.Rep.GetUserByUID(profile.Uid)
	if err != nil {
		return http.StatusNotFound, err
	}

	(*profile).SignForm.Name = base.Name
	(*profile).Phone = base.Phone
	(*profile).Email = base.Email
	(*profile).Uid = base.Uid

	return http.StatusOK, nil
}

func (uc *UserUseCase) UpdateUserPhotos(uid int, newImages *forms.EImageList) models.WorkMessage {
	var imageNames []string
	for index := range *newImages {
		imageNames = append(imageNames, (*newImages)[index].ImgName)
		(*newImages)[index].ImgBase64 = ""
	}
	return uc.Rep.UpdateUserPhotos(uid, &imageNames)
}

func (uc *UserUseCase) UpdateUserBase(form *forms.SignForm) (int, error) {
	usr, err := uc.Rep.GetUserByUID(form.Uid)
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

	if err := uc.Rep.UpdUserGeneral(inf, usr); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (uc *UserUseCase) GetSmallEventsForUser(smallEvents *models.SmallEventList, uid int) models.WorkMessage {
	er := eventRep.NewSqlEventRepository(db.ConnectToDB())
	code, err := er.GetSmallEventsForUser(smallEvents, uid)
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	} else {
		return models.WorkMessage{
			Request: nil,
			Message: "",
			Status:  code,
		}
	}
}

func (uc *UserUseCase) GetUserOwnEvents(ownEvents *models.OwnEventsList, uid int) models.WorkMessage {
	er := eventRep.NewSqlEventRepository(db.ConnectToDB())

	// Get Mid events
	midEventList := models.MidEventList{}
	code, err := er.GetOwnMidEvents(&midEventList, uid)
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}

	// Get Small events
	smallEventList := models.SmallEventList{}
	code, err = er.GetSmallEventsForUser(&smallEventList, uid)
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}

	// Assign and return
	ownEvents.MidEvents = midEventList
	ownEvents.SmallEvents = smallEventList
	return models.WorkMessage{
			Request: nil,
			Message: "",
			Status:  http.StatusOK,
		}
	// TODO: rewrite in goroutines
}
