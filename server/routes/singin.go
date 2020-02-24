package routes

import (
	"../../db"
	"../forms"
	"database/sql"
	"encoding/json"
	htmux "github.com/dimfeld/httptreemux"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.SignUp
	err := decoder.Decode(&form)
	if err != nil {
		GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
		return
	}

	ok := form.Validate()
	if !ok {
		ValidationFailed(w, r)
	}

	user, err := db.GetUserByPhoneOrEmail(db.ConnectToDB(), form.Name, form.Email)
	if err != nil {
		return
	}

	Handle200(w, r, u)
}

func SignInHandler(router *htmux.TreeMux) {
	router.POST("/api/signin", createNewUser)
}
