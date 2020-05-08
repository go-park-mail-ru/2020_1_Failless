package delivery

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/images"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/settings"
	"failless/internal/pkg/user/usecase"
	"log"
	"net/http"

	pb "failless/api/proto/auth"
	json "github.com/mailru/easyjson"
)

////////////// profile part //////////////////

func UpdProfileGeneral(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	var form forms.SignForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	form.Uid = uid
	uc := usecase.GetUseCase()
	if code, err := uc.UpdateUserBase(&form); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, form, http.StatusOK)
}

func UpdUserMetaData(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	var form forms.MetaForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	form.Uid = uid
	uc := usecase.GetUseCase()
	if code, err := uc.UpdateUserMeta(&form); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, form, http.StatusOK)
}

func UpdProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	var form forms.GeneralForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	form.Uid = uid
	if form.Avatar.ImgBase64 != "" {
		if !images.ValidateImage(&form.Avatar, images.Users) {
			network.GenErrorCode(w, r, "image validation failed", http.StatusNotFound)
			return
		}
	}

	uc := usecase.GetUseCase()
	if code, err := uc.UpdateUserInfo(&form); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, form, http.StatusOK)
}

func UploadNewImage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	var form forms.UploadedImage
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	form.Uid = uid
	if form.Uploaded.ImgBase64 == "" ||
		!images.ValidateImage(&form.Uploaded, images.Users) {
		network.GenErrorCode(w, r, "image validation failed", http.StatusNotFound)
		return
	}

	uc := usecase.GetUseCase()
	if err := uc.AddImageToProfile(form.Uid, form.Uploaded.ImgName); err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	network.Jsonify(w, models.WorkMessage{Message: "ok", Status: 200}, http.StatusOK)
}

func GetProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := int64(0)
	if uid = network.GetIdFromRequest(w, r, ps); uid < 0 {
		network.GenErrorCode(w, r, "Uid is incorrect", http.StatusInternalServerError)
		return
	}

	var profile forms.GeneralForm
	profile.Uid = int(uid)
	uc := usecase.GetUseCase()
	if code, err := uc.GetUserInfo(&profile); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, profile, http.StatusOK)
}

func GetSmallEventsForUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := network.GetIdFromRequest(w, r, ps)
	if uid < 0 {
		network.GenErrorCode(w, r, "Uid is incorrect", http.StatusInternalServerError)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	var smallEvents models.SmallEventList
	uc := usecase.GetUseCase()
	message := uc.GetSmallEventsForUser(&smallEvents, int(uid))
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, smallEvents, http.StatusOK)
}

func GetSmallAndMidEventsForUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := network.GetIdFromRequest(w, r, ps)
	if uid < 0 {
		network.GenErrorCode(w, r, "Uid is incorrect", http.StatusInternalServerError)
		return
	}

	var ownEvents models.OwnEventsList
	uc := usecase.GetUseCase()
	message := uc.GetUserOwnEvents(&ownEvents, int(uid))
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, ownEvents, message.Status)
}

func GetProfileSubscriptions(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := network.GetIdFromRequest(w, r, ps)
	if uid < 0 {
		network.GenErrorCode(w, r, "Uid is incorrect", http.StatusInternalServerError)
		return
	}

	var subscriptions models.MidAndBigEventList
	uc := usecase.GetUseCase()
	message := uc.GetUserSubscriptions(&subscriptions, int(uid))
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, subscriptions, message.Status)
}

////////////// user part //////////////////

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	data := r.Context().Value(security.CtxUserKey)
	if data == nil {
		log.Println("data wasn't found")
		network.GenErrorCode(w, r, "User is not authorised", http.StatusUnauthorized)
		return
	}
	network.Jsonify(w, data.(security.UserClaims), http.StatusOK)
}

////////////// authorization part //////////////////

func SignIn(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.CreateLogout(w)

	var form forms.SignForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
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

func Logout(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
	network.CreateLogout(w)
	network.Jsonify(w, models.WorkMessage{Message: "Successfully logout", Status: 200}, 200)
}

////////////// registration part //////////////////

func SignUp(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	data := r.Context().Value(security.CtxUserKey)
	if data != nil {
		network.Jsonify(w, data.(security.UserClaims), http.StatusNotModified)
		return
	}

	var form forms.SignForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	ok := form.Validate()
	if !ok {
		network.ValidationFailed(w, r)
		return
	}

	uc := usecase.GetUseCase()
	user := models.User{
		Phone: form.Phone,
		Email: form.Email,
		Uid:   -1,
	}
	code, err := uc.FillFormIfExist(&user)
	if code != http.StatusNotFound && err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	if user.Uid > 0 {
		network.GenErrorCode(w, r, "User with this information already exist", http.StatusConflict)
		return
	}

	if err := uc.RegisterNewUser(&form); err != nil {
		log.Println("user wasn't registered")
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	form.Password = ""
	network.Jsonify(w, form, http.StatusOK)
}

////////////// feed part //////////////////

func GetUsersFeed(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var searchRequest models.UserRequest
	err := json.UnmarshalFromReader(r.Body, &searchRequest)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	log.Println(searchRequest)

	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}

	// Get FeedUsers to show
	var users []models.UserGeneral
	uc := usecase.GetUseCase()
	if code, err := uc.InitUsersByUserPreferences(&users, &searchRequest); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}
	// Get events and tags about FeedUsers
	var info []forms.GeneralForm
	for i := 0; i < len(users); i++ {
		userForm := forms.GeneralForm{}
		userForm.Uid = (users)[i].Uid
		if code, err := uc.GetUserInfo(&userForm); err != nil {
			network.GenErrorCode(w, r, err.Error(), code)
			return
		}
		info = append(info, userForm)
	}
	feed, err := uc.GetFeedResults(&users, &info)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
	}

	network.Jsonify(w, feed, http.StatusOK)
}
