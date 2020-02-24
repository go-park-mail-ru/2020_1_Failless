package routes

import (
	"../../db"
	"../forms"
	"encoding/json"
	htmux "github.com/dimfeld/httptreemux"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.SignForm
	err := decoder.Decode(&form)
	if err != nil {
		GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
		return
	}

	ok := form.Validate()
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
		err := CreateAuth(w, r, user)
		if err != nil {
			GenErrorCode(w, r, err.Error(), 400)
		}
	}

}

func SignInHandler(router *htmux.TreeMux) {
	router.POST("/api/signin", SignIn)
}
