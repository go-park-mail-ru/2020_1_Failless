package routes

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
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

