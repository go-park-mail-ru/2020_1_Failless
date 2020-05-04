package delivery

import (
	"failless/internal/pkg/event/usecase"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"fmt"
	"log"
	"net/http"

	json "github.com/mailru/easyjson"
)

// Get ALL events ordered by date.
// Deprecated: DO NOT USE IN THE PRODUCTION MODE
func FeedEvents(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uc := usecase.GetUseCase()
	var events models.EventList
	if code, err := uc.InitEventsByTime(&events); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}

func CreateNewEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	var form forms.EventForm
	err := json.UnmarshalFromReader(r.Body, &form)
	if err != nil {
		fmt.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}
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

func CreateSmallEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	var event models.SmallEvent
	err := json.UnmarshalFromReader(r.Body, &event)
	if err != nil {
		fmt.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	err = uc.CreateSmallEvent(&event)
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
	err := json.UnmarshalFromReader(r.Body, &searchRequest)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}
	log.Println(searchRequest)

	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}

	var events models.EventList
	uc := usecase.GetUseCase()
	if code, err := uc.InitEventsByKeyWords(&events, searchRequest.Query, searchRequest.Page); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}

func GetEventsFeed(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var searchRequest models.EventRequest
	err := json.UnmarshalFromReader(r.Body, &searchRequest)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	log.Println(searchRequest)

	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}

	var events models.EventList
	uc := usecase.GetUseCase()
	if code, err := uc.InitEventsByUserPreferences(&events, &searchRequest); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}

func FollowEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if uid := security.CheckCredentials(w, r); uid < 0 {
		return
	}

	var subscription models.EventFollow
	err := json.UnmarshalFromReader(r.Body, &subscription)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	message := uc.FollowEvent(&subscription)

	network.Jsonify(w, message, http.StatusCreated)
}

func GetSearchEvents(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	var searchRequest models.EventRequest
	err := json.UnmarshalFromReader(r.Body, &searchRequest)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	log.Println(searchRequest)

	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}

	var events models.EventResponseList
	uc := usecase.GetUseCase()
	if code, err := uc.SearchEventsByUserPreferences(&events, &searchRequest); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}
