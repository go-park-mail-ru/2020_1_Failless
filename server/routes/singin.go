package routes

import (
	"encoding/json"
	htmux "github.com/dimfeld/httptreemux"
	"github.com/go-park-mailru/2020_1_Failless/db"
	"github.com/go-park-mailru/2020_1_Failless/server/forms"
	"github.com/go-park-mailru/2020_1_Failless/server/utils"
	"log"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("handler work")
	ok, err := utils.IsAuth(w, r)
	if err != nil || ok {
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

	ok = form.Validate()
	if !ok {
		ValidationFailed(w, r)
		return
	}

	user, err := db.GetUserByPhoneOrEmail(db.ConnectToDB(), form.Name, form.Email)
	if err != nil {
		GenErrorCode(w, r, "User doesn't exist", http.StatusNotFound)
		return
	}

	if ComparePasswords(user.Password, form.Password) {
		err := utils.CreateAuth(w, r, user)
		if err != nil {
			GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
			return
		}
	}

}

func SignInHandler(router *htmux.TreeMux) {
	router.POST("/api/signin", SignIn)
}
