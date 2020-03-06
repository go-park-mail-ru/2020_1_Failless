package delivery

import (
	"encoding/json"
	"failless/internal/app/auth/usecase"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/middleware"
	"failless/internal/pkg/network"
	"fmt"
	"log"
	"net/http"
	"os"
)

func SignUp(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("/api/signup")
	data := r.Context().Value(middleware.CtxUserKey)
	if data != nil {
		network.Jsonify(w, data, http.StatusNotModified)
		return
	}

	//r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.SignForm
	err := decoder.Decode(&form)
	if err != nil {
		network.GenErrorCode(w, r, "Invalid Json", http.StatusNotAcceptable)
		return
	}

	log.Println("decoded signup form")
	ok := form.Validate()
	if !ok {
		network.ValidationFailed(w, r)
		network.GenErrorCode(w, r, "Data Error", http.StatusForbidden)
		return
	}
	log.Println("validate signup form")

	user, err := db.GetUserByPhoneOrEmail(db.ConnectToDB(), form.Phone, form.Email)
	if err != nil {
		network.GenErrorCode(w, r, "DB error", http.StatusInternalServerError)
		return
	}

	log.Println(user)
	if user.Uid > 0 {
		network.GenErrorCode(w, r, "User with this information already exist", http.StatusConflict)
		return
	}

	err = usecase.RegisterNewUser(form)
	if err != nil {
		log.Println("user wasn't registered")
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	form.Password = ""
	network.Jsonify(w, form, http.StatusOK)
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
