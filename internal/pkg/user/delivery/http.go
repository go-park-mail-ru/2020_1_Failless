package delivery

import (
	"encoding/json"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/images"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/user/usecase"
	"log"
	"net/http"
)

////////////// profile part //////////////////

func UpdProfileGeneral(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CompareUidsFromURLAndToken(w, r, ps)
	if uid < 0 {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var form forms.SignForm
	err := decoder.Decode(&form)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
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

	decoder := json.NewDecoder(r.Body)
	var form forms.MetaForm
	err := decoder.Decode(&form)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
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

	decoder := json.NewDecoder(r.Body)
	var form forms.GeneralForm
	err := decoder.Decode(&form)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
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

	decoder := json.NewDecoder(r.Body)
	var form forms.UploadedImage
	err := decoder.Decode(&form)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
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

	network.Jsonify(w, network.Message{Message: "ok", Status: 200}, http.StatusOK)
}

func GetProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := 0
	if uid = network.GetIdFromRequest(w, r, &ps); uid < 0 {
		network.GenErrorCode(w, r, "Uid is incorrect", http.StatusInternalServerError)
		return
	}

	var profile forms.GeneralForm
	profile.Uid = uid
	uc := usecase.GetUseCase()
	if code, err := uc.GetUserInfo(&profile); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, profile, http.StatusOK)
}

func GetProfileSubscriptions(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	// TODO: do it in the future
}

////////////// user part //////////////////

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	data := r.Context().Value(security.CtxUserKey)
	if data == nil {
		log.Println("data wasn't found")
		network.GenErrorCode(w, r, "User is not authorised", http.StatusUnauthorized)
		return
	}

	network.Jsonify(w, data, http.StatusOK)
}

////////////// authorization part //////////////////

func SignIn(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.CreateLogout(w)

	decoder := json.NewDecoder(r.Body)
	var form forms.SignForm
	err := decoder.Decode(&form)
	if err != nil {
		network.GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
		return
	}

	if !(form.ValidatePhone() || form.ValidateEmail()) /*|| !(form.ValidatePassword())*/ {
		log.Println("validation error")
		network.ValidationFailed(w, r)
		return
	}

	user := models.User{
		Phone: form.Phone,
		Email: form.Email,
	}
	uc := usecase.GetUseCase()
	if code, err := uc.FillFormIfExist(&user); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	if security.ComparePasswords(user.Password, form.Password) {
		err := network.CreateAuth(w, user)
		if err != nil {
			log.Println("cookie wasn't set")
			network.GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
			return
		}
		log.Println("cookie was set")
	} else {
		network.GenErrorCode(w, r, "Passwords is not equal", http.StatusUnauthorized)
	}

	form.FillFromModel(&user)
	network.Jsonify(w, form, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	network.CreateLogout(w)
	network.Jsonify(w, network.Message{Message: "Successfully logout", Status: 200}, 200)
}

////////////// registration part //////////////////

func SignUp(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	data := r.Context().Value(security.CtxUserKey)
	if data != nil {
		network.Jsonify(w, data, http.StatusNotModified)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var form forms.SignForm
	err := decoder.Decode(&form)
	if err != nil {
		network.GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
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

// debug&test func
func UserDelete(mail string) {
	//err := db.DeleteUser(db.ConnectToDB(), mail)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//log.Println("Success 'UserDelete'")
}

////////////// feed part //////////////////

func GetUsersFeed(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var searchRequest models.UserRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&searchRequest)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
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
		tempinfo := forms.GeneralForm{}
		tempinfo.Uid = users[i].Uid
		if code, err := uc.GetUserInfo(&tempinfo); err != nil {
			network.GenErrorCode(w, r, err.Error(), code)
			return
		}
		info = append(info, tempinfo)
	}

	if len(info) != len(users) {
		network.GenErrorCode(w, r, "invalid lengths", http.StatusBadRequest)
		return
	}

	// Get subscriptions for FeedUsers
	var subs [][]models.Event
	for i := 0; i < len(users); i++ {
		var tempsubs []models.Event
		if code, err := uc.GetUserSubscriptions(&tempsubs, users[i].Uid); err != nil {
			network.GenErrorCode(w, r, err.Error(), code)
			return
		}
		subs = append(subs, tempsubs)
	}

	//Mix up of UserGeneral, GeneralForm and Subs
	type kek struct {
		models.UserGeneral
		Events []models.Event		`json:"events"`
		Tags []models.Tag			`json:"tags"`
		Subs []models.Event			`json:"subscriptions"`
	}

	// Collecting all together
	var result []kek
	for i := 0; i < len(users); i++ {
		tempkek := kek{}
		tempkek.Uid = users[i].Uid
		tempkek.Name = users[i].Name
		tempkek.Photos = users[i].Photos
		tempkek.About = users[i].About
		tempkek.Birthday = users[i].Birthday
		tempkek.Gender = users[i].Gender
		tempkek.Events = info[i].Events
		tempkek.Tags = info[i].Tags
		tempkek.Subs = subs[i]
		result = append(result, tempkek)
	}

	network.Jsonify(w, result, http.StatusOK)
}
