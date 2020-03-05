package middleware
//
//import (
//	"encoding/json"
//	"net/http"
//	"strconv"
//)
//
//type Message struct {
//	Request *http.Request `json:"-"`
//	Message string        `json:"message"`
//	Status  int           `json:"status"`
//}
//
//func GetIdFromRequest(w http.ResponseWriter, r *http.Request, ps *map[string]string) int {
//	uid, err := strconv.Atoi((*ps)["id"])
//	if err != nil {
//		GenErrorCode(w, r, "Incorrect id", http.StatusBadRequest)
//		return -1
//	}
//	return uid
//}
//
//func GenErrorCode(w http.ResponseWriter, r *http.Request, what string, status int) {
//	w.WriteHeader(http.StatusOK)
//	page := Message{r, what, status}
//	output, err := json.Marshal(page)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	w.Header().Set("content-type", "application/json")
//	_, _ = w.Write(output)
//}
//
//func ValidationFailed(w http.ResponseWriter, r *http.Request) {
//	GenErrorCode(w, r, "validation failed", http.StatusBadRequest)
//}

//var AllowedHosts = map[string]struct{}{
//	"http://localhost":           {},
//	"http://localhost:8080":      {},
//	"http://localhost:5000":      {},
//	"http://127.0.0.1":           {},
//	"http://127.0.0.1:8080":      {},
//	"http://127.0.0.1:5000":      {},
//	"https://eventum.rowbot.dev": {},
//}
//
//var AllowedMethods = map[string]struct{}{
//	"GET":     {},
//	"POST":    {},
//	"OPTIONS": {},
//	"HEAD":    {},
//	"PUT":     {},
//}

//func CORS(w http.ResponseWriter, r *http.Request) bool {
//	origin := r.Header.Get("Origin")
//	_, allowed := AllowedHosts[origin]
//	_, allowedMethod := AllowedMethods[r.Method]
//	if allowed && allowedMethod {
//		w.Header().Set("Access-Control-Allow-Origin", origin)
//		w.Header().Set("Access-Control-Allow-Credentials", "true")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//		w.Header().Set("Access-Control-Max-Age", "600")
//		w.Header().Set("Access-Control-Allow-Methods",
//			"GET, POST, OPTIONS, HEAD, PUT")
//	}
//	if r.Method == "OPTIONS" {
//		w.WriteHeader(http.StatusOK)
//		return false
//	}
//	return true
//}
//
//func Jsonify(w http.ResponseWriter, object interface{}, status int) {
//	output, err := json.Marshal(object)
//	if err != nil {
//		http.Error(w, err.Error(), status)
//		return
//	}
//
//	w.Header().Set("content-type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	_, err = w.Write(output)
//	if err != nil {
//		http.Error(w, err.Error(), status)
//		return
//	}
//	log.Println("Sent json")
//}
