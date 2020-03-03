package delivery

import (
	"encoding/json"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"log"
	"net/http"

	htmux "github.com/dimfeld/httptreemux"
)

func SignIn(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.CORS(w, r)
	log.Print("/api/signin")
	_ = security.CreateLogout(w)

	decoder := json.NewDecoder(r.Body)
	var form forms.SignForm
	err := decoder.Decode(&form)
	if err != nil {
		network.GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
		return
	}

	if !(form.ValidatePhone() || form.ValidateEmail()) /*|| !(form.ValidatePassword())*/ {
		log.Println("validation error")
		network.ValidationFailed(w, r)
		return
	}

	user, err := db.GetUserByPhoneOrEmail(db.ConnectToDB(), form.Phone, form.Email)
	if user.Uid < 0 {
		log.Println("user not found")
		network.GenErrorCode(w, r, "User doesn't exist", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("error was occurred")
		log.Println(err.Error())
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	if security.ComparePasswords(user.Password, form.Password) {
		err := security.CreateAuth(w, user)
		if err != nil {
			network.GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		network.GenErrorCode(w, r, "Passwords is not equal", http.StatusUnauthorized)
	}
	form.Password = ""
	form.Name = user.Name
	form.Email = user.Email
	form.Phone = user.Phone
	form.Uid = user.Uid
	network.Jsonify(w, form, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	network.CORS(w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Print("/api/logout")
	_, err := security.IsAuth(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := security.CreateLogout(w); err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
		return
	}

	network.Jsonify(w, network.ErrorMessage{Message: "Successfully logout", Status: 200}, 200)
}

func AuthHandler(router *htmux.TreeMux) {
	router.POST("/api/signin", SignIn)
	router.GET("/api/logout", Logout)
}
