package network

import (
	"failless/internal/pkg/models"
	"log"
	"net/http"

	json "github.com/mailru/easyjson"
)

func Jsonify(w http.ResponseWriter, object json.Marshaler, status int) {
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
	page := models.WorkMessage{Request: r, Message: what, Status: status}
	output, err := json.Marshal(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(output)
}

func ValidationFailed(w http.ResponseWriter, r *http.Request) {
	GenErrorCode(w, r, MessageValidationFailed, http.StatusBadRequest)
}
