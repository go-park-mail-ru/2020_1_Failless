package routes

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_1_Failless/db"
	"github.com/go-park-mail-ru/2020_1_Failless/server/forms"
	"net/http"
	"path"
	"strings"
	"time"
)

type ErrorMessage struct {
	Request *http.Request
	Message string
	status int
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func GenErrorCode(w http.ResponseWriter, r *http.Request, what string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	page := ErrorMessage{r, what, status}
	_ = json.NewEncoder(w).Encode(page)
}


func ValidationFailed(w http.ResponseWriter, r *http.Request) {
	GenErrorCode(w, r, "validation failed", http.StatusBadRequest)
}

func FillProfile(row db.UserInfo) forms.ProfileForm {
	// todo: take pictures from media
	// todo: fill form
	return forms.ProfileForm{
		SignForm: nil,
		Avatar:   forms.EImage{},
		Photos:   nil,
		Gender:   0,
		About:    "",
		Rating:   0,
		Location: db.LocationPoint{},
		Birthday: time.Time{},
	}
}
