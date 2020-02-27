package routes

import (
	"encoding/json"
	"failless/db"
	"failless/server/forms"
	"failless/server/utils"
	htmux "github.com/dimfeld/httptreemux"
	"log"
	"net/http"
)

func UpdProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("/api/profile")
	uid, err := utils.IsAuth(w, r)
	if err != nil || uid > 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.ProfileForm
	err = decoder.Decode(&form)
	if err != nil {
		GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
		return
	}

	if !(form.Validate() && form.ValidateGender()) {
		log.Println("validation error")
		ValidationFailed(w, r)
		return
	}
	if !form.ValidationImage() {
		GenErrorCode(w, r, "image validation failed", http.StatusNotFound)
	}
	var info db.UserInfo
	var user db.User

	if err := form.GetDBFormat(&info, &user); err != nil {
		GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.AddUserInfo(db.ConnectToDB(), user, info); err != nil {
		GenErrorCode(w, r, err.Error(), http.StatusNotFound)
		return
	}
	output, err := json.Marshal(form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(output)
}

func GetProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("/api/profile")
	uid, err := utils.IsAuth(w, r)
	if err != nil || uid < 0 {
		return
	}

	row, err := db.GetProfileInfo(db.ConnectToDB(), uid)
	if err != nil {
		log.Println(err.Error())
		GenErrorCode(w, r, "Profile not found", http.StatusNotFound)
		return
	}
	profile, err := FillProfile(row)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(output)
}

func ProfileHandler(router *htmux.TreeMux) {
	router.POST("/api/profile", UpdProfilePage)
	router.GET("/api/profile", GetProfilePage)
}
