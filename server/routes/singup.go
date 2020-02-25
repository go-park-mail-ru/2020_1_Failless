package routes

import (
	"eventum/db"
	"eventum/server/forms"
	"eventum/server/utils"
	"encoding/json"
	htmux "github.com/dimfeld/httptreemux"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request, ps map[string]string) {
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
	if user.Uid > 0 {
		GenErrorCode(w, r, "User with this information already exist", http.StatusConflict)
		return
	}

	ok = utils.RegisterNewUser(form)
	if !ok {
		GenErrorCode(w, r, "error", 500)
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

func SignUPHandler(router *htmux.TreeMux) {
	router.POST("/api/signup", SignUp)
}
