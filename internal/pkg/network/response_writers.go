package network

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Request *http.Request `json:"-"`
	Message string        `json:"message"`
	Status  int           `json:"status"`
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
	log.Println("Sent json")
}

func GenErrorCode(w http.ResponseWriter, r *http.Request, what string, status int) {
	w.WriteHeader(http.StatusOK)
	page := Message{r, what, status}
	output, err := json.Marshal(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(output)
}

func ValidationFailed(w http.ResponseWriter, r *http.Request) {
	GenErrorCode(w, r, "validation failed", http.StatusBadRequest)
}
