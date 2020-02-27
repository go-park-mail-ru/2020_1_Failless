package routes

import (
	"encoding/json"
	htmux "github.com/dimfeld/httptreemux"
	"github.com/go-park-mail-ru/2020_1_Failless/db"
	"github.com/go-park-mail-ru/2020_1_Failless/server/forms"
	"github.com/go-park-mail-ru/2020_1_Failless/server/utils"
	"log"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	log.Print("/api/signin")
	uid, err := utils.IsAuth(w, r)
	if err != nil || uid > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotModified)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.SignForm
	err = decoder.Decode(&form)
	if err != nil {
		GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
		return
	}

	log.Println(form)

	if !(form.ValidatePhone() || form.ValidateEmail()) /*|| !(form.ValidatePassword())*/ {
		log.Println("validation error")
		ValidationFailed(w, r)
		return
	}
	log.Println("validation passed")

	user, err := db.GetUserByPhoneOrEmail(db.ConnectToDB(), form.Name, form.Email)
	if user.Uid < 0 {
		log.Println("user not found")
		GenErrorCode(w, r, "User doesn't exist", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err.Error())
		GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	if ComparePasswords(user.Password, form.Password) {
		err := utils.CreateAuth(w, user)
		if err != nil {
			GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		GenErrorCode(w, r, "Passwords is not equal", http.StatusUnauthorized)
	}
}


func Logout(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("/api/logout")
	uid, err := utils.IsAuth(w, r)
	if err != nil || uid < 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := utils.CreateLogout(w); err != nil {
		GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

func AuthHandler(router *htmux.TreeMux) {
	router.POST("/api/signin", SignIn)
	router.GET("/api/logout", Logout)
}
