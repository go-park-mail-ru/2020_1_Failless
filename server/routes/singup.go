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

func SignUp(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("/api/signup")
	uid, err := utils.IsAuth(w, r)
	if err != nil || uid > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotModified)
		return
	}
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
		GenErrorCode(w, r, "DB error", http.StatusInternalServerError)
		return
	}
	if user.Uid > 0 {
		GenErrorCode(w, r, "User with this information already exist", http.StatusConflict)
		return
	}

	err = utils.RegisterNewUser(form)
	if err != nil {
		log.Println("user wasn't registered")
		GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("user was registered")
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

func SignUPHandler(router *htmux.TreeMux) {
	router.POST("/api/signup", SignUp)
}
