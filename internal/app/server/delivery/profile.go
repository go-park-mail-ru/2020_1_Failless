package delivery

import (
	"encoding/json"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/network"
	"log"
	"net/http"
	"time"
)

func UpdProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	//if !middleware.CORS(w, r) {
	//	return
	//}
	log.Print("/api/profile")
	//_, err := security.IsAuth(w, r)
	//if err != nil {
	//	middleware.GenErrorCode(w, r, "auth required", http.StatusUnauthorized)
	//	return
	//}
	var key int
	data, ok := r.Context().Value(key).(*network.Claims)
	if data == nil || !ok || data.Uid > 0 {
		network.Jsonify(w, data, http.StatusNotModified)
		return
	}

	uid := 0
	if uid = network.GetIdFromRequest(w, r, &ps); uid < 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.ProfileForm
	err := decoder.Decode(&form)
	if err != nil {
		form1 := forms.ProfileForm{
			SignForm: forms.SignForm{
				Name:     "me",
				Phone:    "88005553535",
				Email:    "rowbot@dev.dev",
				Password: "root12345",
			},
			Avatar: forms.EImage{
				ImgBase64: "stststst",
				ImgName:   "name",
			},
			Photos: []forms.EImage{
				{
					ImgBase64: "stststst",
					ImgName:   "name",
				},
			},
			Gender: 0,
			About:  "about me",
			Rating: 0,
			Location: db.LocationPoint{
				Longitude: 1212.1,
				Latitude:  1212.1,
				Accuracy:  12121,
			},
			Birthday: time.Now(),
		}
		network.Jsonify(w, form1, 200)
		return
	}

	// if !(form.Validate() && form.ValidateGender()) {
	// 	log.Println("validation error")
	// 	ValidationFailed(w, r)
	// 	return
	// }
	if form.Avatar.ImgBase64 != "" {
		if !form.ValidationImage() {
			network.GenErrorCode(w, r, "image validation failed", http.StatusNotFound)
			return
		}
	}
	var info db.UserInfo
	var user db.User

	if err := form.GetDBFormat(&info, &user); err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Uid = uid
	if err := db.AddUserInfo(db.ConnectToDB(), user, info); err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusNotFound)
		return
	}

	form.Avatar.ImgBase64 = ""
	for _, item := range form.Photos {
		item.ImgBase64 = ""
		item.Img = nil
	}
	network.Jsonify(w, form, http.StatusOK)
}

func GetProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	//if !middleware.CORS(w, r) {
	//	return
	//}
	log.Println("/api/profile")
	//_, err := security.IsAuth(w, r)
	//if err != nil {
	//	middleware.GenErrorCode(w, r, "auth required", http.StatusUnauthorized)
	//	return
	//}
	var key int
	data, ok := r.Context().Value(key).(*network.Claims)
	if data == nil || !ok || data.Uid > 0 {
		network.Jsonify(w, data, http.StatusNotModified)
		return
	}

	uid := 0
	if uid = network.GetIdFromRequest(w, r, &ps); uid < 0 {
		return
	}

	log.Println(uid)
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
	log.Println(base.Name)
	profile.SignForm.Name = base.Name
	profile.Password = ""
	network.Jsonify(w, profile, http.StatusOK)
	//output, err := json.Marshal(profile)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//w.Header().Set("content-type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//_, _ = w.Write(output)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	//if !middleware.CORS(w, r) {
	//	return
	//}
	log.Println("/api/getuser")
	//data, err := security.IsAuth(w, r)
	//if err != nil {
	//	middleware.GenErrorCode(w, r, "User is not authorised", http.StatusUnauthorized)
	//	return
	//}

	var key int
	data, ok := r.Context().Value(key).(*network.Claims)
	if data == nil || !ok || data.Uid > 0 {
		network.Jsonify(w, data, http.StatusNotModified)
		return
	}

	// todo: find data for user info
	network.Jsonify(w, data, http.StatusOK)
}

//func ProfileHandler(router *htmux.TreeMux) {
//	router.POST("/api/profile/:id", UpdProfilePage)
//	router.GET("/api/profile/:id", GetProfilePage)
//	router.GET("/api/getuser", GetUserInfo)
//}
