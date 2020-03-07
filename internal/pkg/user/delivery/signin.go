package delivery

import (
	"encoding/json"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/user/usecase"
	"log"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	_ = network.CreateLogout(w)

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

	var user models.User
	uc := usecase.GetUseCase()
	if code, err := uc.FillFormIfExist(&user); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	if security.ComparePasswords(user.Password, form.Password) {
		err := network.CreateAuth(w, user)
		if err != nil {
			network.GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		network.GenErrorCode(w, r, "Passwords is not equal", http.StatusUnauthorized)
	}

	form.FillFromModel(&user)
	network.Jsonify(w, form, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if err := network.CreateLogout(w); err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
		return
	}

	network.Jsonify(w, network.Message{Message: "Successfully logout", Status: 200}, 200)
}
