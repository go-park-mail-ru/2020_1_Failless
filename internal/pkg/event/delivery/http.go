package delivery

import (
	"encoding/json"
	"failless/internal/pkg/event/usecase"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/middleware"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"log"
	"net/http"
)

// Get ALL events ordered by date.
// DO NOT USE IN THE PRODUCTION MODE
func FeedEvents(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uc := usecase.GetUseCase()
	var events []models.Event
	if code, err := uc.InitEventsByTime(&events); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}

func CreateNewEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	data := r.Context().Value(middleware.CtxUserKey)
	if data == nil {
		network.GenErrorCode(w, r, "auth required", http.StatusUnauthorized)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var form forms.EventForm
	err := decoder.Decode(&form)
	if !form.Validate() {
		// TODO: add error code from error code table
		network.GenErrorCode(w, r, "incorrect data", http.StatusBadRequest)
		return
	}
	uc := usecase.GetUseCase()
	event, err := uc.CreateEvent(form)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	network.Jsonify(w, event, http.StatusOK)
}

// Get events limited by number strings with offset from JSON (POST parameter)
// Limit have to be set in the /configs/*/settings.go file using global variable
// UseCaseConf
func GetEventsByKeyWords(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var searchRequest models.EventRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&searchRequest)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
		return
	}
	log.Println(searchRequest)

	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}

	var events []models.Event
	uc := usecase.GetUseCase()
	if code, err := uc.InitEventsByKeyWords(&events, searchRequest.Query, searchRequest.Page); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}

