package routes

import (
	"encoding/json"
	"failless/db"
	"failless/server/forms"
	"failless/server/utils"
	"fmt"
	htmux "github.com/dimfeld/httptreemux"
	"log"
	"net/http"
	"os"
)

func SignUp(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("/api/signup")
	if !CORS(w, r) {
		return
	}

	data, err := utils.IsAuth(w, r)
	if data.Uid > 0 {
		Jsonify(w, data, http.StatusNotModified)
		//GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
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

	log.Println("decoded signup form")
	ok := form.Validate()
	if !ok {
		ValidationFailed(w, r)
		GenErrorCode(w, r, "Data Error", http.StatusForbidden)
		return
	}
	log.Println("validate signup form")

	user, err := db.GetUserByPhoneOrEmail(db.ConnectToDB(), form.Phone, form.Email)
	if err != nil {
		GenErrorCode(w, r, "DB error", http.StatusInternalServerError)
		return
	}

	log.Println(user)
	if user.Uid > 0 {
		GenErrorCode(w, r, "User with this information already exist", http.StatusConflict)
		return
	}

	err = forms.RegisterNewUser(form)
	if err != nil {
		log.Println("user wasn't registered")
		GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	form.Password = ""
	output, err := json.Marshal(form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(output)
}

func OptionsReq(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if !CORS(w, r) {
		return
	}
}

// debug&test func
func UserDelete(mail string) {
	err := db.DeleteUser(db.ConnectToDB(), mail)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Success 'UserDelete'")
}

func SignUPHandler(router *htmux.TreeMux) {
	router.POST("/api/signup", SignUp)
	router.OptionsHandler = OptionsReq
}
