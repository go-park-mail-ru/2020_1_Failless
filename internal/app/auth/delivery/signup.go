package delivery

import (
	"encoding/json"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/middleware"
	"failless/internal/pkg/security"
	"fmt"
	"log"
	"net/http"
	"os"

	htmux "github.com/dimfeld/httptreemux"
)

func SignUp(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("/api/signup")
	if !middleware.CORS(w, r) {
		return
	}

	data, err := security.IsAuth(w, r)
	if data.Uid > 0 {
		middleware.Jsonify(w, data, http.StatusNotModified)
		//GenErrorCode(w, r, err.Error(), middleware.StatusUnauthorized)
		return
	}
	r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.SignForm
	err = decoder.Decode(&form)
	if err != nil {
		middleware.GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
		return
	}

	log.Println("decoded signup form")
	ok := form.Validate()
	if !ok {
		middleware.ValidationFailed(w, r)
		middleware.GenErrorCode(w, r, "Data Error", http.StatusForbidden)
		return
	}
	log.Println("validate signup form")

	user, err := db.GetUserByPhoneOrEmail(db.ConnectToDB(), form.Phone, form.Email)
	if err != nil {
		middleware.GenErrorCode(w, r, "DB error", http.StatusInternalServerError)
		return
	}

	log.Println(user)
	if user.Uid > 0 {
		middleware.GenErrorCode(w, r, "User with this information already exist", http.StatusConflict)
		return
	}

	err = forms.RegisterNewUser(form)
	if err != nil {
		log.Println("user wasn't registered")
		middleware.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
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
	if !middleware.CORS(w, r) {
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
