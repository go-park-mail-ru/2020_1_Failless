package delivery

import (
	"encoding/json"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/middleware"
	"failless/internal/pkg/network"
	"failless/internal/pkg/user/usecase"
	"log"
	"net/http"
)

func UpdProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	var key int
	data, ok := r.Context().Value(key).(*network.Claims)
	if data == nil || !ok || data.Uid > 0 {
		network.GenErrorCode(w, r, "auth required", http.StatusUnauthorized)
		return
	}

	uid := 0
	if uid = network.GetIdFromRequest(w, r, &ps); uid < 0 {
		network.GenErrorCode(w, r, "Uid is incorrect", http.StatusInternalServerError)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.ProfileForm
	err := decoder.Decode(&form)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
		return
	}

	form.Uid = uid
	if form.Avatar.ImgBase64 != "" {
		if !form.ValidationImage() {
			network.GenErrorCode(w, r, "image validation failed", http.StatusNotFound)
			return
		}
	}
	uc := usecase.GetUseCase()
	if err, code := uc.UpdateUserInfo(&form); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}
	//var info db.UserInfo
	//var user db.User
	//
	//if err := form.GetDBFormat(&info, &user); err != nil {
	//	network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//user.Uid = uid
	//if err := db.AddUserInfo(db.ConnectToDB(), user, info); err != nil {
	//	network.GenErrorCode(w, r, err.Error(), http.StatusNotFound)
	//	return
	//}
	//
	//form.Avatar.ImgBase64 = ""
	//for _, item := range form.Photos {
	//	item.ImgBase64 = ""
	//	item.Img = nil
	//}
	network.Jsonify(w, form, http.StatusOK)
}

func GetProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("GET: /api/profile")
	uid := 0
	if uid = network.GetIdFromRequest(w, r, &ps); uid < 0 {
		network.GenErrorCode(w, r, "Uid is incorrect", http.StatusInternalServerError)
		return
	}

	row, err := db.GetProfileInfo(db.ConnectToDB(), uid)
	if err != nil {
		log.Println(err.Error())
		network.GenErrorCode(w, r, "Profile not found", http.StatusNotFound)
		return
	}

	var profile forms.ProfileForm
	err = profile.FillProfile(row)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	base, err := db.GetUserByUID(db.ConnectToDB(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	profile.SignForm.Name = base.Name
	profile.Password = ""
	log.Println(profile)
	network.Jsonify(w, profile, http.StatusOK)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("/api/getuser")
	data := r.Context().Value(middleware.CtxUserKey)
	log.Println(data)
	if data == nil {
		log.Println("data wasn't found")
		network.GenErrorCode(w, r, "User is not authorised", http.StatusUnauthorized)
		return
	}
	network.Jsonify(w, data, http.StatusOK)
}
