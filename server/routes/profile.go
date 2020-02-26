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

func UpdProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("handler work")
	uid, err := utils.IsAuth(w, r)
	if err != nil || uid > 0 {
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

	// todo: added images download and update db data
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


}

func ProfileHandler(router *htmux.TreeMux) {
	router.POST("/api/profile", UpdProfilePage)
	router.GET("/api/profile", GetProfilePage)
}
