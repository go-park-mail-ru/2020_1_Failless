package delivery

import (
	"failless/internal/pkg/event/usecase"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/images"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
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

// Get events limited by number strings with offset from JSON (POST parameter)
// UserCount have to be set in the /configs/*/settings.go file using global variable
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

func OLDGetEventsFeed(w http.ResponseWriter, r *http.Request, ps map[string]string) {
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

	var events models.MidAndBigEventList
	uc := usecase.GetUseCase()
	if code, err := uc.SearchEventsByUserPreferences(&events, &searchRequest); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}
	network.Jsonify(w, events, http.StatusOK)
}

func GetSmallEvents(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}


	uc := usecase.GetUseCase()
	events, err := uc.GetSmallEventsByUID(int64(uid))
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}

func CreateSmallEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	// Get form
	r.Header.Set("Content-Type", "application/json")
	var eventForm forms.SmallEventForm
	err := json.UnmarshalFromReader(r.Body, &eventForm)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	// Check form
	if !eventForm.Validate() {
		network.GenErrorCode(w, r, "Small event validation failed. Check server logs", http.StatusBadRequest)
		return
	}

	// Upload pics
	for iii := 0; iii < len(eventForm.Photos); iii++ {
		if eventForm.Photos[iii].ImgBase64 == "" ||
			!images.ValidateImage(&eventForm.Photos[iii], images.Events) {
			network.GenErrorCode(w, r, "Image validation failed. Check server logs", http.StatusBadRequest)
		}
	}

	uc := usecase.GetUseCase()
	event, err := uc.CreateSmallEvent(&eventForm)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	network.Jsonify(w, event, http.StatusOK)
}

func UpdateSmallEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	var event models.SmallEvent
	err := json.UnmarshalFromReader(r.Body, &event)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	code, err := uc.UpdateSmallEvent(&event)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, event, code)
}

func DeleteSmallEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	eid := network.GetEIdFromRequest(w, r, ps)
	if eid < 0 {
		network.GenErrorCode(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	uc := usecase.GetUseCase()
	message := uc.DeleteSmallEvent(uid, eid)

	network.Jsonify(w, message, message.Status)
}

func CreateMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	var midEventForm forms.MidEventForm
	err := json.UnmarshalFromReader(r.Body, &midEventForm)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	if !midEventForm.Validate() {
		network.GenErrorCode(w, r, "incorrect data", http.StatusBadRequest)
		return
	}

	// Upload pics
	for iii := 0; iii < len(midEventForm.Photos); iii++ {
		if midEventForm.Photos[iii].ImgBase64 == "" ||
			!images.ValidateImage(&midEventForm.Photos[iii], images.Events) {
			network.GenErrorCode(w, r, "Image validation failed. Check server logs", http.StatusBadRequest)
		}
	}

	uc := usecase.GetUseCase()
	midEvent, message := uc.CreateMidEvent(&midEventForm)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, midEvent, message.Status)
}

func GetMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	network.GenErrorCode(w, r, "Not implemented", 502)
}

func UpdateMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	network.GenErrorCode(w, r, "Not implemented", 502)

}

func DeleteMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	network.GenErrorCode(w, r, "Not implemented", 502)
}

func JoinMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if uid := security.CheckCredentials(w, r); uid < 0 {
		return
	}

	if eid := network.GetEIdFromRequest(w, r, ps); eid < 0 {
		network.GenErrorCode(w, r, "Error in retrieving eid from url", http.StatusBadRequest)
	}

	var subscription models.EventFollow
	err := json.UnmarshalFromReader(r.Body, &subscription)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	message := uc.JoinMidEvent(&subscription)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, message, message.Status)
}

func LeaveMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if uid := security.CheckCredentials(w, r); uid < 0 {
		return
	}

	if eid := network.GetEIdFromRequest(w, r, ps); eid < 0 {
		network.GenErrorCode(w, r, "Error in retrieving eid from url", http.StatusBadRequest)
	}

	var subscription models.EventFollow
	err := json.UnmarshalFromReader(r.Body, &subscription)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	message := uc.LeaveMidEvent(&subscription)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, message, message.Status)
}

func CreateBigEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	panic("impement me!")
}

func GetBigEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	panic("impement me!")
}

func UpdateBigEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	panic("impement me!")
}

func DeleteBigEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	panic("impement me!")
}

func AddVisitorForBigEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	panic("impement me!")
}

func RemoveVisitorForBigEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	panic("impement me!")
}
