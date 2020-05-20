package delivery

//go:generate mockgen -destination=../mocks/mock_delivery.go -package=mocks failless/internal/pkg/user Delivery

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/images"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/settings"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/usecase"
	"log"
	"net/http"

	pb "failless/api/proto/auth"
	json "github.com/mailru/easyjson"
)

type userDelivery struct {
	UseCase user.UseCase
}

func GetDelivery() user.Delivery {
	return &userDelivery{
		UseCase: usecase.GetUseCase(),
	}
}

////////////// profile part //////////////////

func (ud *userDelivery) UpdProfileGeneral(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	var form forms.SignForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	form.Uid = uid
	if code, err := ud.UseCase.UpdateUserBase(&form); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, form, http.StatusOK)
}

func (ud *userDelivery) UpdUserAbout(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	var about models.UserAbout
	err := json.UnmarshalFromReader(r.Body, &about)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	message := ud.UseCase.UpdateUserAbout(uid, about.About)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, about, message.Status)
}

func (ud *userDelivery) UpdUserTags(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	var tags models.UserTags
	err := json.UnmarshalFromReader(r.Body, &tags)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	message := ud.UseCase.UpdateUserTags(uid, tags.Tags)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, tags, message.Status)
}

func (ud *userDelivery) UpdUserPhotos(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	var newImages forms.EImageList
	err := json.UnmarshalFromReader(r.Body, &newImages)
	if err != nil {
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	for index := range newImages {
		if newImages[index].ImgBase64 != "" {
			if newImages[index].ImgName == "" || !images.ValidateImage(&newImages[index], images.Users) {
				network.GenErrorCode(w, r, images.MessageImageValidationFailed, http.StatusNotFound)
				return
			}
		}
	}

	message := ud.UseCase.UpdateUserPhotos(uid, &newImages)
	if message.Status != http.StatusOK {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, newImages, message.Status)
}

func (ud *userDelivery) GetProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := int64(0)
	if uid = network.GetIdFromRequest(w, r, ps); uid < 0 {
		return
	}

	var profile forms.GeneralForm
	profile.Uid = int(uid)
	if code, err := ud.UseCase.GetUserInfo(&profile); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, profile, http.StatusOK)
}

func (ud *userDelivery) GetSmallEventsForUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	// TODO: insecure
	uid := network.GetIdFromRequest(w, r, ps)
	if uid < 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	var smallEvents models.SmallEventList
	message := ud.UseCase.GetSmallEventsForUser(&smallEvents, int(uid))
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, smallEvents, http.StatusOK)
}

func (ud *userDelivery) GetSmallAndMidEventsForUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := network.GetIdFromRequest(w, r, ps)
	if uid < 0 {
		return
	}

	var ownEvents models.OwnEventsList
	message := ud.UseCase.GetUserOwnEvents(&ownEvents, int(uid))
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, ownEvents, message.Status)
}

func (ud *userDelivery) GetProfileSubscriptions(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := network.GetIdFromRequest(w, r, ps)
	if uid < 0 {
		return
	}

	var subscriptions models.MidAndBigEventList
	message := ud.UseCase.GetUserSubscriptions(&subscriptions, int(uid))
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, subscriptions, message.Status)
}

////////////// user part //////////////////

func (ud *userDelivery) GetUserInfo(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	data := r.Context().Value(security.CtxUserKey)
	if data == nil {
		network.GenErrorCode(w, r, network.MessageErrorAuthRequired, http.StatusUnauthorized)
		return
	}
	network.Jsonify(w, data.(security.UserClaims), http.StatusOK)
}

////////////// authorization part //////////////////

func (ud *userDelivery) SignIn(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.CreateLogout(w)

	var form forms.SignForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	if !(form.ValidatePhone() || form.ValidateEmail()) /*|| !(form.ValidatePassword())*/ {
		log.Println("validation error")
		network.ValidationFailed(w, r)
		return
	}

	authReply, err := settings.AuthClient.Authorize(r.Context(), &pb.AuthRequest{
		Phone:    form.Phone,
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	if !authReply.Ok {
		network.GenErrorCode(w, r, authReply.Message, http.StatusUnauthorized)
		return
	}

	token, err := settings.AuthClient.GetToken(r.Context(), &pb.AuthRequest{
		Uid:      authReply.Cred.Uid,
		Phone:    form.Phone,
		Email:    form.Email,
		Password: form.Password,
	})
	if token != nil {
		network.CreateAuthMS(&w, token.Token)
	}

	form.FillFromAuthReply(authReply)
	network.Jsonify(w, form, http.StatusOK)
}

func (ud *userDelivery) Logout(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
	network.CreateLogout(w)
	network.Jsonify(w, models.WorkMessage{Message: network.MessageSuccessfulLogout, Status: http.StatusOK}, http.StatusOK)
}

////////////// registration part //////////////////

func (ud *userDelivery) SignUp(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	data := r.Context().Value(security.CtxUserKey)
	if data != nil {
		network.Jsonify(w, data.(security.UserClaims), http.StatusNotModified)
		return
	}

	var form forms.SignForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	ok := form.Validate()
	if !ok {
		network.ValidationFailed(w, r)
		return
	}

	user := models.User{
		Phone: form.Phone,
		Email: form.Email,
		Uid:   -1,
	}
	code, err := ud.UseCase.FillFormIfExist(&user)
	if code != http.StatusNotFound && err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	if user.Uid > 0 {
		network.GenErrorCode(w, r, "User with this information already exist", http.StatusConflict)
		return
	}

	if err := ud.UseCase.RegisterNewUser(&form); err != nil {
		log.Println("user wasn't registered")
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	form.Password = ""
	network.Jsonify(w, form, http.StatusOK)
}

////////////// feed part //////////////////

func (ud *userDelivery) GetUsersFeed(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var searchRequest models.UserRequest
	err := json.UnmarshalFromReader(r.Body, &searchRequest)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	log.Println(searchRequest)

	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}

	// Get FeedUsers to show
	var users []models.UserGeneral
	if code, err := ud.UseCase.InitUsersByUserPreferences(&users, &searchRequest); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	feedResults, message := ud.UseCase.GetFeedResultsFor(uid, &users)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
	}

	network.Jsonify(w, feedResults, http.StatusOK)
}
