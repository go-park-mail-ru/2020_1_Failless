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
	status  int
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

var AllowedHosts = map[string]struct{}{
	"http://localhost":           {},
	"http://localhost:8080":      {},
	"http://localhost:5000":      {},
	"http://127.0.0.1":           {},
	"http://127.0.0.1:8080":      {},
	"http://127.0.0.1:5000":      {},
	"https://eventum.rowbot.dev": {},
}

var AllowedMethods = map[string]struct{}{
	"GET":     {},
	"POST":    {},
	"OPTIONS": {},
	"HEAD":    {},
	"PUT":     {},
}

func CORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	_, allowed := AllowedHosts[origin]
	_, allowedMethod := AllowedHosts[r.Method]
	if allowed && allowedMethod {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "600")
		w.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, OPTIONS, HEAD, PUT")
	}
}

func Jsonify(w http.ResponseWriter, object interface{}, status int) {
	output, err := json.Marshal(object)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
}